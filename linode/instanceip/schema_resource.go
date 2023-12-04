package instanceip

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceIPVPCNAT1To1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the VPC interface mapped to this IP.",
			},
			"subnet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the VPC subnet the corresponding interface is attached to.",
			},
			"vpc_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the VPC the corresponding interface is attached to.",
			},
		},
	}
}

var resourceSchema = map[string]*schema.Schema{
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Linode to allocate an IPv4 address for.",
		Required:    true,
		ForceNew:    true,
	},
	"public": {
		Type:        schema.TypeBool,
		Description: "Whether the IPv4 address is public or private.",
		Default:     true,
		Optional:    true,
		ForceNew:    true,
	},

	"address": {
		Type:        schema.TypeString,
		Description: "The resulting IPv4 address.",
		Computed:    true,
	},
	"gateway": {
		Type:        schema.TypeString,
		Description: "The default gateway for this address",
		Computed:    true,
	},
	"prefix": {
		Type:        schema.TypeInt,
		Description: "The number of bits set in the subnet mask.",
		Computed:    true,
	},
	"rdns": {
		Type:        schema.TypeString,
		Description: "The reverse DNS assigned to this address.",
		Optional:    true,
		Computed:    true,
	},
	"region": {
		Type:        schema.TypeString,
		Description: "The region this IP resides in.",
		Computed:    true,
	},
	"subnet_mask": {
		Type:        schema.TypeString,
		Description: "The mask that separates host bits from network bits for this address.",
		Computed:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The type of IP address.",
		Computed:    true,
	},
	"vpc_nat_1_1": {
		Type:        schema.TypeList,
		Description: "Contains information about the NAT 1:1 mapping of a public IP address to a VPC subnet.",
		Computed:    true,
		Elem:        resourceIPVPCNAT1To1(),
	},

	"apply_immediately": {
		Type: schema.TypeBool,
		Description: "If true, the instance will be rebooted to update network interfaces. " +
			"This functionality is not affected by the `skip_implicit_reboots` provider argument.",
		Optional: true,
	},
}
