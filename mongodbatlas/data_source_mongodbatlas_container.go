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
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"atlas_cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioned": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
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
	d.Set("region", c.RegionName)
	d.Set("vpc_id", c.VpcID)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	return nil
}
