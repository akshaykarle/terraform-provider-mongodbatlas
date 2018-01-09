package mongodbatlas

import (
	"fmt"
	"log"

	"github.com/akshaykarle/mongodb-atlas-go/mongodb"
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
	return nil
}
