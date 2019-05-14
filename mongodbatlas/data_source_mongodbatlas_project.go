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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_count": {
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
	if err := d.Set("org_id", p.OrgID); err != nil {
		return fmt.Errorf("error setting org_id for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("name", p.Name); err != nil {
		return fmt.Errorf("error setting name for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("created", p.Created); err != nil {
		return fmt.Errorf("error setting created for resource %s: %s", d.Id(), err)
	}
	if err := d.Set("cluster_count", p.ClusterCount); err != nil {
		return fmt.Errorf("error setting cluster_count for resource %s: %s", d.Id(), err)
	}

	return nil
}
