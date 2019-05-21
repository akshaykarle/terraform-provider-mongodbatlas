package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerCreate,
		Read:   resourceContainerRead,
		Update: resourceContainerUpdate,
		Delete: resourceContainerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceContainerImportState,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"atlas_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gcp_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provisioned": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"private_ip_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.Container{
		AtlasCidrBlock: d.Get("atlas_cidr_block").(string),
		ProviderName:   d.Get("provider_name").(string),
	}

	if params.ProviderName == "AWS" {
		params.RegionName = d.Get("region").(string)
	}

	if params.ProviderName == "GCP" {
		params.GcpProjectID = d.Get("gcp_project_id").(string)
		params.NetworkName = d.Get("network_name").(string)
	}

	container, _, err := client.Containers.Create(d.Get("group").(string), &params)

	if err != nil {
		return fmt.Errorf("Error creating MongoDB Container: %s", err)
	}
	d.SetId(container.ID)
	log.Printf("[INFO] MongoDB Container ID: %s", d.Id())

	if d.Get("private_ip_mode").(bool) {
		log.Printf("[INFO] Attempting to enable PrivateIPMode")
		_, err := client.PrivateIPMode.Enable(d.Get("group").(string))

		if err != nil {
			return fmt.Errorf("Error attempting to enable PrivateIPMode: %s", err)
		}
	}

	return resourceContainerRead(d, meta)
}

func resourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	c, _, err := client.Containers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
	}

	if err := d.Set("atlas_cidr_block", c.AtlasCidrBlock); err != nil {
		log.Printf("[WARN] Error setting atlas_cidr_block for (%s): %s", d.Id(), err)
	}
	if err := d.Set("provider_name", c.ProviderName); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", d.Id(), err)
	}
	if err := d.Set("identifier", c.ID); err != nil {
		log.Printf("[WARN] Error setting identifier for (%s): %s", d.Id(), err)
	}
	if err := d.Set("provisioned", c.Provisioned); err != nil {
		log.Printf("[WARN] Error setting provisioned for (%s): %s", d.Id(), err)
	}

	if d.Get("provider_name").(string) == "AWS" {
		if err := d.Set("region", c.RegionName); err != nil {
			log.Printf("[WARN] Error setting region for (%s): %s", d.Id(), err)
		}
		if err := d.Set("vpc_id", c.VpcID); err != nil {
			log.Printf("[WARN] Error setting vpc_id for (%s): %s", d.Id(), err)
		}
	}
	if d.Get("provider_name").(string) == "GCP" {
		if err := d.Set("gcp_project_id", c.GcpProjectID); err != nil {
			log.Printf("[WARN] Error setting gcp_project_id for (%s): %s", d.Id(), err)
		}
		if err := d.Set("network_name", c.NetworkName); err != nil {
			log.Printf("[WARN] Error setting network_name for (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.Containers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
	}

	if d.HasChange("atlas_cidr_block") {
		c.AtlasCidrBlock = d.Get("atlas_cidr_block").(string)
		requestUpdate = true
	}
	if d.HasChange("provider_name") {
		c.ProviderName = d.Get("provider_name").(string)
		requestUpdate = true
	}
	if d.HasChange("region") {
		c.RegionName = d.Get("region").(string)
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
	if d.HasChange("private_ip_mode") {
		if d.Get("private_ip_mode").(bool) {
			_, err := client.PrivateIPMode.Enable(d.Get("group").(string))
			if err != nil {
				return fmt.Errorf("Error enabling PrivateIPMode on MongoDB Container %s: %s", d.Id(), err)
			}
		} else {
			_, err := client.PrivateIPMode.Disable(d.Get("group").(string))
			if err != nil {
				return fmt.Errorf("Error disabling PrivateIPMode on MongoDB Container %s: %s", d.Id(), err)
			}
		}
	}

	if requestUpdate {
		c.ID = ""
		_, _, err := client.Containers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
		}
	}

	return resourceContainerRead(d, meta)
}

func resourceContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	group := d.Get("group").(string)

	_, err := client.Containers.Delete(group, d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting MongoDB Container %s: %s", d.Id(), err)
	}

	if d.Get("private_ip_mode").(bool) {
		_, err = client.PrivateIPMode.Disable(group)
		if err != nil {
			return fmt.Errorf("Error enabling PrivateIPMode on MongoDB Container %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceContainerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a container, use the format {group id}-{container id}")
	}
	gid := parts[0]
	containerID := parts[1]
	client := meta.(*ma.Client)
	c, _, err := client.Containers.Get(gid, containerID)
	if err != nil {
		return nil, fmt.Errorf("Error reading MongoDB Container %s: %s", containerID, err)
	}

	d.SetId(c.ID)
	if err := d.Set("atlas_cidr_block", c.AtlasCidrBlock); err != nil {
		log.Printf("[WARN] Error setting atlas_cidr_block for (%s): %s", d.Id(), err)
	}
	if err := d.Set("provider_name", c.ProviderName); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", d.Id(), err)
	}
	if err := d.Set("identifier", c.ID); err != nil {
		log.Printf("[WARN] Error setting identifier for (%s): %s", d.Id(), err)
	}
	if err := d.Set("provisioned", c.Provisioned); err != nil {
		log.Printf("[WARN] Error setting provisioned for (%s): %s", d.Id(), err)
	}

	if d.Get("provider_name").(string) == "AWS" {
		if err := d.Set("region", c.RegionName); err != nil {
			log.Printf("[WARN] Error setting region for (%s): %s", d.Id(), err)
		}
		if err := d.Set("vpc_id", c.VpcID); err != nil {
			log.Printf("[WARN] Error setting vpc_id for (%s): %s", d.Id(), err)
		}
	}
	if d.Get("provider_name").(string) == "GCP" {
		if err := d.Set("gcp_project_id", c.GcpProjectID); err != nil {
			log.Printf("[WARN] Error setting gcp_project_id for (%s): %s", d.Id(), err)
		}
		if err := d.Set("network_name", c.NetworkName); err != nil {
			log.Printf("[WARN] Error setting network_name for (%s): %s", d.Id(), err)
		}
	}

	return []*schema.ResourceData{d}, nil

}
