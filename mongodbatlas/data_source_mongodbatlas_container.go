package mongodbatlas

import (
	"fmt"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceContainer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceContainerRead,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	id := d.Get("identifier").(string)
	group := d.Get("group").(string)

	c, _, err := client.Containers.Get(group, id)
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container with Id %s in Group %s: %s", id, group, err)
	}

	d.SetId(c.ID)
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

	return nil
}
