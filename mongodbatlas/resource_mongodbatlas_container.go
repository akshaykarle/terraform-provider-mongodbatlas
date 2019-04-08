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
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"atlas_cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gcp_project_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provisioned": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"private_ip_mode": &schema.Schema{
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

	// Supply needfuls when GCP
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
		log.Printf("Attempting to enable PrivateIPMode")
		_, err := client.PrivateIPMode.EnablePrivateIPMode(d.Get("group").(string))

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

	d.Set("atlas_cidr_block", c.AtlasCidrBlock)
	d.Set("provider_name", c.ProviderName)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	if c.ProviderName == "GCP" {
		d.Set("gcp_project_id", c.GcpProjectID)
		d.Set("network_name", c.NetworkName)
	} else if c.ProviderName == "AWS" {
		d.Set("vpc_id", c.VpcID)
		d.Set("region", c.RegionName)
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
		if d.Get("private_ip_mode").(bool) == false {
			_, err := client.PrivateIPMode.DisablePrivateIPMode(d.Get("group").(string))
			if err != nil {
				return fmt.Errorf("Error disabling PrivateIPMode on MongoDB Container %s: %s", d.Id(), err)
			}
		} else if d.Get("private_ip_mode").(bool) == true {
			_, err := client.PrivateIPMode.EnablePrivateIPMode(d.Get("group").(string))
			if err != nil {
				return fmt.Errorf("Error enabling PrivateIPMode on MongoDB Container %s: %s", d.Id(), err)
			}
		}
	}

	if requestUpdate {
		_, _, err := client.Containers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
		}
	}

	return resourceContainerRead(d, meta)
}

func resourceContainerDelete(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("group", gid)
	d.Set("atlas_cidr_block", c.AtlasCidrBlock)
	d.Set("provider_name", c.ProviderName)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	if d.Get("provider_name").(string) == "AWS" {
		d.Set("region", c.RegionName)
		d.Set("vpc_id", c.VpcID)
	}
	if d.Get("provider_name").(string) == "GCP" {
		d.Set("gcp_project_id", c.GcpProjectID)
		d.Set("network_name", c.NetworkName)
	}

	return []*schema.ResourceData{d}, nil

}
