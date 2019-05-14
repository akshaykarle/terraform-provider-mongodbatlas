package mongodbatlas

import (
	"fmt"
	"log"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertConfigurationCreate,
		Read:   resourceAlertConfigurationRead,
		Update: resourceAlertConfigurationUpdate,
		Delete: resourceAlertConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"event_type_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"matchers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"notifications": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"interval_min": {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(key, oldValue, newValue string, d *schema.ResourceData) bool {
								return oldValue == "2147483647"
							},
						},
						"delay_min": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"email_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"sms_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"team_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mobile_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"notification_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"room_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"channel_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"org_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flow_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flowdock_api_token": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"victor_ops_api_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"victor_ops_routing_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ops_genie_api_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"metric_threshold": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "AVERAGE",
						},
						"units": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlertConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := alertConfigurationFromResourceData(d)

	alert, _, err := client.AlertConfigurations.Create(d.Get("group").(string), params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Alert Configuration: %s", err)
	}
	d.SetId(alert.ID)
	log.Printf("[INFO] MongoDB Alert Configuration ID: %s", d.Id())

	resourceDataFromAlertConfiguration(alert, d)
	return nil
}

func resourceAlertConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	alert, response, err := client.AlertConfigurations.Get(d.Get("group").(string), d.Id())
	if err != nil {
		if response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading MongoDB Alert Configuration %s: %s", d.Id(), err)
	}

	resourceDataFromAlertConfiguration(alert, d)
	return nil
}

func resourceAlertConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := alertConfigurationFromResourceData(d)

	alert, _, err := client.AlertConfigurations.Update(d.Get("group").(string), d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error updating MongoDB Alert Configuration: %s", err)
	}

	resourceDataFromAlertConfiguration(alert, d)
	return nil
}

func resourceAlertConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB Alert Configuration destroy: %v", d.Id())

	response, err := client.AlertConfigurations.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		if response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting MongoDB Alert Configuration: %s", err)
	}

	return nil
}

func alertConfigurationFromResourceData(d *schema.ResourceData) *ma.AlertConfiguration {
	alertConfiguration := ma.AlertConfiguration{
		Enabled:       d.Get("enabled").(bool),
		EventTypeName: d.Get("event_type_name").(string),
		Notifications: readNotificationsFromSchema(d.Get("notifications").([]interface{})),
		Matchers:      readMatchersFromSchema(d.Get("matchers").([]interface{})),
	}
	if v, ok := d.GetOk("metric_threshold"); ok {
		alertConfiguration.MetricThreshold = readMetricThresholdFromSchema(v.([]interface{})[0].(map[string]interface{}))
	}
	return &alertConfiguration
}

func resourceDataFromAlertConfiguration(alertConfiguration *ma.AlertConfiguration, d *schema.ResourceData) {
	notifications := []interface{}{}
	for _, notification := range alertConfiguration.Notifications {
		note := map[string]interface{}{
			"type_name":              notification.TypeName,
			"interval_min":           notification.IntervalMin,
			"delay_min":              notification.DelayMin,
			"email_enabled":          notification.EmailEnabled,
			"sms_enabled":            notification.SMSEnabled,
			"username":               notification.Username,
			"team_id":                notification.TeamID,
			"email_address":          notification.EmailAddress,
			"mobile_number":          notification.MobileNumber,
			"notification_token":     notification.NotificationToken,
			"room_name":              notification.RoomName,
			"channel_name":           notification.ChannelName,
			"api_token":              notification.APIToken,
			"org_name":               notification.OrgName,
			"flow_name":              notification.FlowName,
			"flowdock_api_token":     notification.FlowdockAPIToken,
			"service_key":            notification.ServiceKey,
			"victor_ops_api_key":     notification.VictorOpsAPIKey,
			"victor_ops_routing_key": notification.VictorOpsRoutingKey,
			"ops_genie_api_key":      notification.OpsGenieAPIKey,
		}
		notifications = append(notifications, note)
	}
	if err := d.Set("notifications", notifications); err != nil {
		log.Printf("[WARN] Error setting notifications for (%s): %s", d.Id(), err)
	}

	matchers := []interface{}{}
	for _, matcher := range alertConfiguration.Matchers {
		match := map[string]interface{}{
			"field_name": matcher.FieldName,
			"operator":   matcher.Operator,
			"value":      matcher.Value,
		}
		matchers = append(matchers, match)
	}
	if err := d.Set("matchers", matchers); err != nil {
		log.Printf("[WARN] Error setting matchers for (%s): %s", d.Id(), err)
	}

	metricThreshold := map[string]interface{}{
		"metric_name": alertConfiguration.MetricThreshold.MetricName,
		"operator":    alertConfiguration.MetricThreshold.Operator,
		"threshold":   alertConfiguration.MetricThreshold.Threshold,
		"units":       alertConfiguration.MetricThreshold.Units,
		"mode":        alertConfiguration.MetricThreshold.Mode,
	}

	if err := d.Set("metric_threshold", metricThreshold); err != nil {
		log.Printf("[WARN] Error setting metric threshold for (%s): %s", d.Id(), err)
	}

	if err := d.Set("event_type_name", alertConfiguration.EventTypeName); err != nil {
		log.Printf("[WARN] Error setting event_type_name for (%s): %s", d.Id(), err)
	}
	if err := d.Set("enabled", alertConfiguration.Enabled); err != nil {
		log.Printf("[WARN] Error setting enabled for (%s): %s", d.Id(), err)
	}
	if err := d.Set("group", alertConfiguration.GroupID); err != nil {
		log.Printf("[WARN] Error setting group for (%s): %s", d.Id(), err)
	}
}

func readMetricThresholdFromSchema(thresholdMap map[string]interface{}) (threshold ma.MetricThreshold) {
	fmt.Println(thresholdMap)
	threshold = ma.MetricThreshold{
		MetricName: thresholdMap["metric_name"].(string),
		Operator:   thresholdMap["operator"].(string),
		Threshold:  thresholdMap["threshold"].(float64),
		Units:      thresholdMap["units"].(string),
		Mode:       thresholdMap["mode"].(string),
	}
	return threshold
}

func readMatchersFromSchema(matchersMap []interface{}) (matchers []ma.Matcher) {
	matchers = make([]ma.Matcher, len(matchersMap))
	for i, r := range matchersMap {
		matcherMap := r.(map[string]interface{})

		matchers[i] = ma.Matcher{
			FieldName: matcherMap["field_name"].(string),
			Operator:  matcherMap["operator"].(string),
			Value:     matcherMap["value"].(string),
		}
	}
	return matchers
}

func readNotificationsFromSchema(notificationsMap []interface{}) (notifications []ma.Notification) {
	notifications = make([]ma.Notification, len(notificationsMap))
	for i, r := range notificationsMap {
		notificationMap := r.(map[string]interface{})

		notifications[i] = ma.Notification{
			TypeName:            notificationMap["type_name"].(string),
			IntervalMin:         notificationMap["interval_min"].(int),
			DelayMin:            notificationMap["delay_min"].(int),
			EmailEnabled:        notificationMap["email_enabled"].(bool),
			SMSEnabled:          notificationMap["sms_enabled"].(bool),
			Username:            notificationMap["username"].(string),
			TeamID:              notificationMap["team_id"].(string),
			EmailAddress:        notificationMap["email_address"].(string),
			MobileNumber:        notificationMap["mobile_number"].(string),
			NotificationToken:   notificationMap["notification_token"].(string),
			RoomName:            notificationMap["room_name"].(string),
			ChannelName:         notificationMap["channel_name"].(string),
			APIToken:            notificationMap["api_token"].(string),
			OrgName:             notificationMap["org_name"].(string),
			FlowName:            notificationMap["flow_name"].(string),
			FlowdockAPIToken:    notificationMap["flowdock_api_token"].(string),
			ServiceKey:          notificationMap["service_key"].(string),
			VictorOpsAPIKey:     notificationMap["victor_ops_api_key"].(string),
			VictorOpsRoutingKey: notificationMap["victor_ops_routing_key"].(string),
			OpsGenieAPIKey:      notificationMap["ops_genie_api_key"].(string),
		}
	}
	return notifications
}
