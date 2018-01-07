package mongodbatlas

import (
	"fmt"
	"log"

	"github.com/akshaykarle/mongodb-atlas-go/mongodb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mongodb_major_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backing_provider": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_size_gb": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"replication_factor": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	providerSettings := mongodb.ProviderSettings{
		ProviderName:        d.Get("provider_name").(string),
		BackingProviderName: d.Get("backing_provider").(string),
		RegionName:          d.Get("region").(string),
		InstanceSizeName:    d.Get("size").(string),
	}
	params := mongodb.Cluster{
		Name:                d.Get("name").(string),
		MongoDBMajorVersion: d.Get("mongodb_major_version").(string),
		ProviderSettings:    providerSettings,
		BackupEnabled:       d.Get("backup").(bool),
		ReplicationFactor:   d.Get("replication_factor").(int),
		DiskSizeGB:          d.Get("disk_size_gb").(float64),
	}

	cluster, _, err := client.Clusters.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Cluster: %s", err)
	}
	d.SetId(cluster.ID)
	log.Printf("[INFO] MongoDB Cluster ID: %s", d.Id())

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)
	c, _, err := client.Clusters.Get(d.Get("group").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Cluster %s: %s", d.Get("name").(string), err)
	}

	d.Set("name", c.Name)
	d.Set("group", c.GroupID)
	d.Set("mongodb_major_version", c.MongoDBMajorVersion)
	d.Set("backup", c.BackupEnabled)
	d.Set("size", c.ProviderSettings.InstanceSizeName)
	d.Set("provider_name", c.ProviderSettings.ProviderName)
	d.Set("backing_provider", c.ProviderSettings.BackingProviderName)
	d.Set("region", c.ProviderSettings.RegionName)
	d.Set("disk_size_gb", c.DiskSizeGB)
	d.Set("replication_factor", c.ReplicationFactor)
	d.Set("identifier", c.ID)
	d.Set("state", c.StateName)

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
