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
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"atlas_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	id := d.Get("identifier").(string)
	group := d.Get("group").(string)

	c, _, err := client.Containers.Get(group, id)
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container with Id %s in Group %s: %s", id, group, err)
	}

	d.SetId(c.ID)
	if err := d.Set("atlas_cidr_block", c.AtlasCidrBlock); err != nil {
		return fmt.Errorf("error setting atlas_cidr_block for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("provider_name", c.ProviderName); err != nil {
		return fmt.Errorf("error setting provider_name for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("identifier", c.ID); err != nil {
		return fmt.Errorf("error setting identifier for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("provisioned", c.Provisioned); err != nil {
		return fmt.Errorf("error setting provisioned for resource %s: %s", d.Id(), err)
	}

	if d.Get("provider_name").(string) == "AWS" {
		if err := d.Set("region", c.RegionName); err != nil {
			return fmt.Errorf("error setting region for resource %s: %s", d.Id(), err)
		}
		if err := d.Set("vpc_id", c.VpcID); err != nil {
			return fmt.Errorf("error setting vpc_id for resource %s: %s", d.Id(), err)
		}
	}
	if d.Get("provider_name").(string) == "GCP" {
		if err := d.Set("gcp_project_id", c.GcpProjectID); err != nil {
			return fmt.Errorf("error setting gcp_project_id for resource %s: %s", d.Id(), err)
		}
		if err := d.Set("network_name", c.NetworkName); err != nil {
			return fmt.Errorf("error setting network_name for resource %s: %s", d.Id(), err)
		}
	}

	return nil
}
