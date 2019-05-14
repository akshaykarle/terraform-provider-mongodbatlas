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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	if err := d.Set("org_id", p.OrgID); err != nil {
		log.Printf("[WARN] Error setting org_id for (%s): %s", d.Id(), err)
	}
	if err := d.Set("name", p.Name); err != nil {
		log.Printf("[WARN] Error setting name for (%s): %s", d.Id(), err)
	}
	if err := d.Set("created", p.Created); err != nil {
		log.Printf("[WARN] Error setting created for (%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_count", p.ClusterCount); err != nil {
		log.Printf("[WARN] Error setting cluster_count for (%s): %s", d.Id(), err)
	}

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
