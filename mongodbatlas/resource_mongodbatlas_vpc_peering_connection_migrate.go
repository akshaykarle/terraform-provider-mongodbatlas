package mongodbatlas

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVpcPeeringConnectionResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"error_state_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVpcPeeringConnectionStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if v, ok := rawState["provider_name"]; (ok && v == "") || !ok {
		rawState["provider_name"] = "AWS"
	}

	return rawState, nil
}
