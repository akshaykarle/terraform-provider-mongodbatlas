package mongodbatlas

import (
	"fmt"
	"log"

	"github.com/akshaykarle/mongodb-atlas-go/mongodb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerCreate,
		Read:   resourceContainerRead,
		Update: resourceContainerUpdate,
		Delete: resourceContainerDelete,

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
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
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
		},
	}
}

func resourceContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	params := mongodb.Container{
		AtlasCidrBlock: d.Get("atlas_cidr_block").(string),
		ProviderName:   d.Get("provider_name").(string),
		RegionName:     d.Get("region").(string),
	}

	container, _, err := client.Containers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Container: %s", err)
	}
	d.SetId(container.ID)
	log.Printf("[INFO] MongoDB Container ID: %s", d.Id())

	return resourceContainerRead(d, meta)
}

func resourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	c, _, err := client.Containers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
	}

	d.Set("atlas_cidr_block", c.AtlasCidrBlock)
	d.Set("provider_name", c.ProviderName)
	d.Set("region", c.RegionName)
	d.Set("vpc_id", c.VpcID)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceContainerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
