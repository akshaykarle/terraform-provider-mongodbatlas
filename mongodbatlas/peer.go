package mongodbatlas

import (
	"fmt"
	"log"
	"time"

	"github.com/akshaykarle/mongodb-atlas-go/mongodb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePeer() *schema.Resource {
	return &schema.Resource{
		Create: resourcePeerCreate,
		Read:   resourcePeerRead,
		Update: resourcePeerUpdate,
		Delete: resourcePeerDelete,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"error_state_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	params := mongodb.Peer{
		RouteTableCidrBlock: d.Get("route_table_cidr_block").(string),
		VpcID:               d.Get("vpc_id").(string),
		AwsAccountID:        d.Get("aws_account_id").(string),
		ContainerID:         d.Get("container_id").(string),
	}

	peer, _, err := client.Peers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error initiating MongoDB Peering connection: %s", err)
	}
	d.SetId(peer.ID)
	log.Printf("[INFO] MongoDB Peering ID: %s", d.Id())

	log.Println("[INFO] Waiting for MongoDB VPC Peering Connection to be available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourcePeerStateRefreshFunc(d.Id(), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourcePeerRead(d, meta)
}

func resourcePeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	p, _, err := client.Peers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
	}

	d.Set("route_table_cidr_block", p.RouteTableCidrBlock)
	d.Set("vpc_id", p.VpcID)
	d.Set("aws_account_id", p.AwsAccountID)
	d.Set("identifier", p.ID)
	d.Set("container_id", p.ContainerID)
	d.Set("connection_id", p.ConnectionID)
	d.Set("status_name", p.StatusName)
	d.Set("error_state_name", p.ErrorStateName)

	return nil
}

func resourcePeerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mongodb.Client)

	log.Printf("[DEBUG] MongoDB VPC Peering connection destroy: %v", d.Id())
	_, err := client.Peers.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB VPC Peering connection %s: %s", d.Id(), err)
	}

	log.Println("[INFO] Waiting for MongoDB VPC Peering Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE", "PENDING_ACCEPTANCE", "INITIATING", "FINALIZING", "TERMINATING"},
		Target:     []string{"DELETED"},
		Refresh:    resourcePeerStateRefreshFunc(d.Id(), d.Get("group").(string), client),
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

func resourcePeerStateRefreshFunc(id, group string, client *mongodb.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.Peers.Get(group, id)
		if err != nil {
			if resp.StatusCode == 404 {
				return 42, "DELETED", nil
			}
			log.Printf("Error reading MongoDB VPC Peering connection %s: %s", id, err)
			return nil, "", err
		}

		if p.StatusName != "" {
			log.Printf("[DEBUG] MongoDB Peer status for cluster: %s: %s", id, p.StatusName)
		}

		return p, p.StatusName, nil
	}
}
