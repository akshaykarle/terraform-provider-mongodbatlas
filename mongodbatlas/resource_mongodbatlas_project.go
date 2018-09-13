package mongodbatlas

import (
	"fmt"
	"log"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.Project{
		OrgID: d.Get("org_id").(string),
		Name:  d.Get("name").(string),
	}

	p, _, err := client.Projects.Create(&params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Project IP Projects: %s", err)
	}
	d.SetId(p.ID)
	log.Printf("[INFO] MongoDB Project ID: %s", d.Id())

	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	p, resp, err := client.Projects.Get(d.Id())
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading MongoDB Project %s: %s", d.Id(), err)
	}

	d.Set("org_id", p.OrgID)
	d.Set("name", p.Name)
	d.Set("created", p.Created)
	d.Set("cluster_count", p.ClusterCount)

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB Project destroy: %v", d.Id())
	_, err := client.Projects.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB Project %s: %s", d.Id(), err)
	}

	return nil
}
