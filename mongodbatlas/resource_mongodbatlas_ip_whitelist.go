package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var createMutex = &sync.Mutex{}

func resourceIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPWhitelistCreate,
		Read:   resourceIPWhitelistRead,
		Update: resourceIPWhitelistUpdate,
		Delete: resourceIPWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIPWhiteListImportState,
		},

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_address"},
			},
			"ip_address": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cidr_block"},
			},
			"comment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 80),
			},
		},
	}
}

func resourceIPWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	createMutex.Lock()
	defer createMutex.Unlock()

	client := meta.(*ma.Client)
	cidrBlock := d.Get("cidr_block").(string)
	ip := d.Get("ip_address").(string)

	if cidrBlock != "" && ip != "" {
		// cidrBlock & ip are mutually exclusive, use cidrBlock if both are set
		ip = ""
	}

	params := []ma.Whitelist{
		ma.Whitelist{
			CidrBlock: cidrBlock,
			GroupID:   d.Get("group").(string),
			IPAddress: ip,
			Comment:   d.Get("comment").(string),
		},
	}

	log.Printf("[DEBUG] Creating MongoDB Project IP Whitelist with CIDR block: %v and IP Address: %v", cidrBlock, ip)
	whitelists, _, err := client.Whitelist.Create(d.Get("group").(string), params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Project IP Whitelist: %s", err)
	}
	for _, w := range whitelists {
		if (cidrBlock != "" && w.CidrBlock == cidrBlock) || (ip != "" && w.IPAddress == ip) {
			d.SetId(w.CidrBlock)
			log.Printf("[INFO] MongoDB Project IP Whitelist ID: %s", d.Id())

			return resourceIPWhitelistRead(d, meta)
		}
	}
	return fmt.Errorf("MongoDB Project IP Whitelist with CIDR block: %s and IP Address: %s could not be found in the response from MongoDB Atlas", cidrBlock, ip)
}

func resourceIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	w, _, err := client.Whitelist.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Project IP Whitelist %s: %s", d.Id(), err)
	}

	d.Set("cidr_block", w.CidrBlock)
	d.Set("ip_address", w.IPAddress)
	d.Set("group", w.GroupID)
	d.Set("comment", w.Comment)

	return nil
}

func resourceIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIPWhitelistCreate(d, meta)
}

func resourceIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB Project IP Whitelist destroy: %v", d.Id())
	_, err := client.Whitelist.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB Project IP Whitelist %s: %s", d.Id(), err)
	}

	return nil
}

func resourceIPWhiteListImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ma.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import an ip whitelist, use the format {group id}-{cidr block}")
	}
	gid := parts[0]
	cidr := parts[1]

	ip, _, err := client.Whitelist.Get(gid, cidr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import ip whitelist %s in group %s, error: %s", cidr, gid, err.Error())
	}

	d.SetId(ip.CidrBlock)
	d.Set("group", ip.GroupID)

	return []*schema.ResourceData{d}, nil
}
