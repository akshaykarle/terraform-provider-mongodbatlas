package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceClusterImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},

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
			},
			"backup": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
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
				ForceNew: false,
			},
			"disk_size_gb": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"replication_factor": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"num_shards": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"paused": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disk_gb_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"mongodb_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mongo_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mongo_uri_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mongo_uri_with_options": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"replication_spec": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"electable_nodes": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"read_only_nodes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	providerSettings := ma.ProviderSettings{
		ProviderName:        d.Get("provider_name").(string),
		BackingProviderName: d.Get("backing_provider").(string),
		RegionName:          d.Get("region").(string),
		InstanceSizeName:    d.Get("size").(string),
	}
	autoScaling := ma.AutoScaling{
		DiskGBEnabled: d.Get("disk_gb_enabled").(bool),
	}
	params := ma.Cluster{
		Name:                d.Get("name").(string),
		MongoDBMajorVersion: d.Get("mongodb_major_version").(string),
		ProviderSettings:    providerSettings,
		BackupEnabled:       d.Get("backup").(bool),
		ReplicationFactor:   d.Get("replication_factor").(int),
		ReplicationSpec:     readReplicationSpecsFromSchema(d.Get("replication_spec").(*schema.Set).List()),
		DiskSizeGB:          d.Get("disk_size_gb").(float64),
		NumShards:           d.Get("num_shards").(int),
		Paused:              d.Get("paused").(bool),
		AutoScaling:         autoScaling,
	}

	cluster, _, err := client.Clusters.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Cluster: %s", err)
	}
	d.SetId(cluster.ID)
	log.Printf("[INFO] MongoDB Cluster ID: %s", d.Id())

	log.Println("[INFO] Waiting for MongoDB Cluster to be available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterStateRefreshFunc(d.Get("name").(string), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	c, resp, err := client.Clusters.Get(d.Get("group").(string), d.Get("name").(string))
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading MongoDB Cluster %s: %s", d.Get("name").(string), err)
	}

	replicationSpecs := []interface{}{}
	for region, replicationSpec := range c.ReplicationSpec {
		spec := map[string]interface{}{
			"region":          region,
			"priority":        replicationSpec.Priority,
			"electable_nodes": replicationSpec.ElectableNodes,
			"read_only_nodes": replicationSpec.ReadOnlyNodes,
		}
		replicationSpecs = append(replicationSpecs, spec)
	}

	if err := d.Set("replication_spec", replicationSpecs); err != nil {
		log.Printf("[WARN] Error setting replication specs set for (%s): %s", d.Get("name"), err)
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
	d.Set("disk_gb_enabled", c.AutoScaling.DiskGBEnabled)
	d.Set("replication_factor", c.ReplicationFactor)
	d.Set("identifier", c.ID)
	d.Set("state", c.StateName)
	d.Set("num_shards", c.NumShards)
	d.Set("paused", c.Paused)
	d.Set("mongodb_version", c.MongoDBVersion)
	d.Set("mongo_uri", c.MongoURI)
	d.Set("mongo_uri_updated", c.MongoURIUpdated)
	d.Set("mongo_uri_with_options", c.MongoURIWithOptions)

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.Clusters.Get(d.Get("group").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Cluster %s: %s", d.Get("name").(string), err)
	}

	if d.HasChange("backup") {
		c.BackupEnabled = d.Get("backup").(bool)
		requestUpdate = true
	}
	if d.HasChange("size") {
		c.ProviderSettings.InstanceSizeName = d.Get("size").(string)
		requestUpdate = true
	}
	if d.HasChange("disk_size_gb") {
		c.DiskSizeGB = d.Get("disk_size_gb").(float64)
		requestUpdate = true
	}
	if d.HasChange("replication_factor") {
		c.ReplicationFactor = d.Get("replication_factor").(int)
		requestUpdate = true
	}
	if d.HasChange("replication_spec") {
		c.ReplicationSpec = readReplicationSpecsFromSchema(d.Get("replication_spec").(*schema.Set).List())
		requestUpdate = true
	}
	if d.HasChange("num_shards") {
		c.NumShards = d.Get("num_shards").(int)
		requestUpdate = true
	}
	if d.HasChange("paused") {
		c.Paused = d.Get("paused").(bool)
		requestUpdate = true
	}
	if d.HasChange("disk_gb_enabled") {
		c.AutoScaling.DiskGBEnabled = d.Get("disk_gb_enabled").(bool)
		requestUpdate = true
	}

	if requestUpdate {
		// Set read-only fields to an empty string to make the API happy
		c.StateName = ""
		c.MongoDBVersion = ""
		c.MongoURI = ""
		c.MongoURIWithOptions = ""
		c.MongoURIUpdated = ""
		_, _, err := client.Clusters.Update(d.Get("group").(string), d.Get("name").(string), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Cluster %s: %s", d.Get("name").(string), err)
		}

		log.Println("[INFO] Waiting for MongoDB Cluster to be updated")

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
			Target:     []string{"IDLE"},
			Refresh:    resourceClusterStateRefreshFunc(d.Get("name").(string), d.Get("group").(string), client),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
			Delay:      30 * time.Second, // Wait 30 secs before starting
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}
	}
	return resourceClusterRead(d, meta)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB Cluster destroy: %v", d.Id())
	_, err := client.Clusters.Delete(d.Get("group").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB Cluster %s: %s", d.Get("name").(string), err)
	}

	log.Println("[INFO] Waiting for MongoDB Cluster to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterStateRefreshFunc(d.Get("name").(string), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return nil
}

func resourceClusterImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ma.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a cluster, use the format {group id}-{cluster name}")
	}
	gid := parts[0]
	name := parts[1]

	c, _, err := client.Clusters.Get(gid, name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import cluster %s in group %s, error: %s", name, gid, err.Error())
	}

	d.SetId(c.ID)
	d.Set("name", c.Name)
	d.Set("group", c.GroupID)

	return []*schema.ResourceData{d}, nil
}

func resourceClusterStateRefreshFunc(name, group string, client *ma.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.Clusters.Get(group, name)
		if err != nil {
			if resp.StatusCode == 404 {
				return 42, "DELETED", nil
			}
			log.Printf("Error reading MongoDB Cluster %s: %s", name, err)
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] MongoDB Cluster status for cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func readReplicationSpecsFromSchema(replicationSpecs []interface{}) map[string]ma.ReplicationSpec {
	specs := map[string]ma.ReplicationSpec{}
	for _, r := range replicationSpecs {
		replicationSpec := r.(map[string]interface{})
		specs[replicationSpec["region"].(string)] = ma.ReplicationSpec{
			Priority:       replicationSpec["priority"].(int),
			ElectableNodes: replicationSpec["electable_nodes"].(int),
			ReadOnlyNodes:  replicationSpec["read_only_nodes"].(int),
		}
	}
	return specs
}
