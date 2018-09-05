package mongodbatlas

import (
	"fmt"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	name := d.Get("name").(string)
	p, _, err := client.Projects.GetByName(name)
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Project with name %s: %s", name, err)
	}

	d.SetId(p.ID)
	d.Set("org_id", p.OrgID)
	d.Set("name", p.Name)
	d.Set("created", p.Created)
	d.Set("cluster_count", p.ClusterCount)

	return nil
}
