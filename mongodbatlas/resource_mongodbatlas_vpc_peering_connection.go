package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcPeeringConnectionCreate,
		Read:   resourceVpcPeeringConnectionRead,
		Update: resourceVpcPeeringConnectionUpdate,
		Delete: resourceVpcPeeringConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVpcPeeringConnectionImportState,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceVpcPeeringConnectionResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceVpcPeeringConnectionStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gcp_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"error_state_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVpcPeeringConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.Peer{
		ContainerID:  d.Get("container_id").(string),
		ProviderName: d.Get("provider_name").(string),
	}

	if d.Get("provider_name").(string) == "AWS" {
		params.RouteTableCidrBlock = d.Get("route_table_cidr_block").(string)
		params.VpcID = d.Get("vpc_id").(string)
		params.AwsAccountID = d.Get("aws_account_id").(string)
	}
	if d.Get("provider_name").(string) == "GCP" {
		params.GcpProjectID = d.Get("gcp_project_id").(string)
		params.NetworkName = d.Get("network_name").(string)
	}

	peer, _, err := client.Peers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error initiating MongoDB Peering connection: %s", err)
	}
	d.SetId(peer.ID)
	log.Printf("[INFO] MongoDB Peering ID: %s", d.Id())

	log.Println("[INFO] Waiting for MongoDB VPC Peering Connection to be available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE", "WAITING_FOR_USER"},
		Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      60 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceVpcPeeringConnectionRead(d, meta)
}

func resourceVpcPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	p, _, err := client.Peers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
	}

	if d.Get("provider_name").(string) == "AWS" {
		if err := d.Set("route_table_cidr_block", p.RouteTableCidrBlock); err != nil {
			log.Printf("[WARN] Error settingroute_table_cidr_block for (%s): %s", d.Id(), err)
		}
		if err := d.Set("vpc_id", p.VpcID); err != nil {
			log.Printf("[WARN] Error setting vpc_id for (%s): %s", d.Id(), err)
		}
		if err := d.Set("aws_account_id", p.AwsAccountID); err != nil {
			log.Printf("[WARN] Error setting aws_account_id for (%s): %s", d.Id(), err)
		}
		if err := d.Set("status_name", p.StatusName); err != nil {
			log.Printf("[WARN] Error setting status_name for (%s): %s", d.Id(), err)
		}
		if err := d.Set("error_state_name", p.ErrorStateName); err != nil {
			log.Printf("[WARN] Error setting error_state_name for (%s): %s", d.Id(), err)
		}
		if err := d.Set("connection_id", p.ConnectionID); err != nil {
			log.Printf("[WARN] Error setting connection_id for (%s): %s", d.Id(), err)
		}
	}
	if d.Get("provider_name").(string) == "GCP" {
		if err := d.Set("network_name", p.NetworkName); err != nil {
			log.Printf("[WARN] Error setting network_name for (%s): %s", d.Id(), err)
		}
		if err := d.Set("gcp_project_id", p.GcpProjectID); err != nil {
			log.Printf("[WARN] Error setting gcp_project_id for (%s): %s", d.Id(), err)
		}
		if err := d.Set("status", p.Status); err != nil {
			log.Printf("[WARN] Error setting status for (%s): %s", d.Id(), err)
		}
		if err := d.Set("error_message", p.ErrorMessage); err != nil {
			log.Printf("[WARN] Error setting error_message for (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("identifier", p.ID); err != nil {
		log.Printf("[WARN] Error setting identifier for (%s): %s", d.Id(), err)
	}
	if err := d.Set("container_id", p.ContainerID); err != nil {
		log.Printf("[WARN] Error setting container_id for (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceVpcPeeringConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.Peers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
	}

	if d.HasChange("route_table_cidr_block") {
		c.RouteTableCidrBlock = d.Get("route_table_cidr_block").(string)
		requestUpdate = true
	}
	if d.HasChange("aws_account_id") {
		c.AwsAccountID = d.Get("aws_account_id").(string)
		requestUpdate = true
	}
	if d.HasChange("vpc_id") {
		c.VpcID = d.Get("vpc_id").(string)
		requestUpdate = true
	}
	if d.HasChange("gcp_project_id") {
		c.GcpProjectID = d.Get("gcp_project_id").(string)
		requestUpdate = true
	}
	if d.HasChange("network_name") {
		c.NetworkName = d.Get("network_name").(string)
		requestUpdate = true
	}

	if requestUpdate {
		_, _, err := client.Peers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER"},
			Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE", "WAITING_FOR_USER", "AVAILABLE"},
			Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
			Delay:      30 * time.Second, // Wait 30 secs before starting
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}
	}

	return resourceVpcPeeringConnectionRead(d, meta)
}

func resourceVpcPeeringConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB VPC Peering connection destroy: %v", d.Id())
	_, err := client.Peers.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB VPC Peering connection %s: %s", d.Id(), err)
	}

	log.Println("[INFO] Waiting for MongoDB VPC Peering Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE", "PENDING_ACCEPTANCE", "INITIATING", "FINALIZING", "TERMINATING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 10 * time.Second,
		Delay:      60 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return nil
}

func getConnection(client *ma.Client, gid string, connectionID string) (*ma.Peer, error) {
	peer, _, err := client.Peers.Get(gid, connectionID)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import vpc peering %s in group %s, error: %s", connectionID, gid, err.Error())
	}
	return peer, nil
}

func resourceVpcPeeringConnectionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a VPC peering, use the format {group id}-{peering connection id}")
	}
	gid := parts[0]
	connectionID := parts[1]
	client := meta.(*ma.Client)
	peer, err := getConnection(client, gid, connectionID)
	if err != nil {
		return nil, err
	}

	// https://docs.atlas.mongodb.com/reference/api/vpc-get-connection/#example-response
	// Atlas API does not return ProviderName, so we have to guess it from other parameters
	if peer.AwsAccountID != "" {
		if err := d.Set("provider_name", "AWS"); err != nil {
			return nil, fmt.Errorf("Error setting provider name: %v", err)
		}
	} else if peer.GcpProjectID != "" {
		if err := d.Set("provider_name", "GCP"); err != nil {
			return nil, fmt.Errorf("Error setting provider name: %v", err)
		}
	}

	d.SetId(peer.ID)
	if err := d.Set("group", gid); err != nil {
		log.Printf("[WARN] Error setting group for (%s): %s", d.Id(), err)
	}

	return []*schema.ResourceData{d}, nil

}

func resourceVpcPeeringConnectionStateRefreshFunc(id, group string, client *ma.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.Peers.Get(group, id)
		status := ""

		if len(p.StatusName) > 0 {
			status = p.StatusName
		} else if len(p.Status) > 0 {
			status = p.Status
		}

		log.Printf("[INFO] Current status: %s", status)

		if err != nil {
			if resp.StatusCode == 404 {
				return 42, "DELETED", nil
			}
			log.Printf("Error reading MongoDB VPC Peering connection %s: %s", id, err)
			return nil, "", err
		}

		if status != "" {
			log.Printf("[DEBUG] MongoDB Peer status for cluster: %s: %s", id, status)
		}

		return p, status, nil
	}
}
