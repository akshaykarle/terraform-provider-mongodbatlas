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

func resourceVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcPeeringConnectionCreate,
		Read:   resourceVpcPeeringConnectionRead,
		Update: resourceVpcPeeringConnectionUpdate,
		Delete: resourceVpcPeeringConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVpcPeeringConnectionImportState,
		},

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

func resourceVpcPeeringConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.Peer{
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
		Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceVpcPeeringConnectionRead(d, meta)
}

func resourceVpcPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

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

func resourceVpcPeeringConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.Peers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
	}

	if d.HasChange("route_table_cidr_block") {
		c.RouteTableCidrBlock = d.Get("route_table_cidr_block").(string)
		requestUpdate = true
	}
	if d.HasChange("aws_account_id") {
		c.AwsAccountID = d.Get("aws_account_id").(string)
		requestUpdate = true
	}
	if d.HasChange("vpc_id") {
		c.VpcID = d.Get("vpc_id").(string)
		requestUpdate = true
	}

	if requestUpdate {
		_, _, err := client.Peers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Peering connection %s: %s", d.Id(), err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"INITIATING", "FINALIZING"},
			Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
			Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
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

	return resourceVpcPeeringConnectionRead(d, meta)
}

func resourceVpcPeeringConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB VPC Peering connection destroy: %v", d.Id())
	_, err := client.Peers.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB VPC Peering connection %s: %s", d.Id(), err)
	}

	log.Println("[INFO] Waiting for MongoDB VPC Peering Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE", "PENDING_ACCEPTANCE", "INITIATING", "FINALIZING", "TERMINATING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceVpcPeeringConnectionStateRefreshFunc(d.Id(), d.Get("group").(string), client),
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

func getConnection(client *ma.Client, gid string, connection_id string) (*ma.Peer, error) {
	peers, _, err := client.Peers.List(gid)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import vpc peering %s in group %s, error: %s", connection_id, gid, err.Error())
	}
	for i := range peers {
		if peers[i].ConnectionID == connection_id {
			return &peers[i], nil
		}
	}
	return nil, fmt.Errorf("Couldn't find vpc peering %s in group %s", connection_id, gid)

}

func resourceVpcPeeringConnectionImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a VPC peering, use the format {group id}-{peering connection id}")
	}
	gid := parts[0]
	connection_id := parts[1]
	client := meta.(*ma.Client)
	peer, err := getConnection(client, gid, connection_id)
	if err != nil {
		return nil, err
	}

	d.SetId(peer.ID)
	d.Set("group", gid)

	return []*schema.ResourceData{d}, nil

}

func resourceVpcPeeringConnectionStateRefreshFunc(id, group string, client *ma.Client) resource.StateRefreshFunc {
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
