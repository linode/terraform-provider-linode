package instancenetworking

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var networkSchema = map[string]*schema.Schema{
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
}

var dataSourceSchema = map[string]*schema.Schema{
	"linode_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Linode for network info.",
		Required:    true,
		ForceNew:    true,
	},
	"ipv4": {
		Type:        schema.TypeList,
		Description: "Information about this Linode’s IPv4 addresses.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"private": {
					Type:        schema.TypeList,
					Description: "A list of private IP Address objects belonging to this Linode.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
				"public": {
					Type:        schema.TypeList,
					Description: "A list of public IP Address objects belonging to this Linode.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
				"reserved": {
					Type:        schema.TypeList,
					Description: "A list of reserved IP Address objects belonging to this Linode.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
				"shared": {
					Type:        schema.TypeList,
					Description: "A list of shared IP Address objects assigned to this Linode.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
			},
		},
	},
	"ipv6": {
		Type:        schema.TypeList,
		Description: "Information about this Linode’s IPv6 addresses.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"global": {
					Type:        schema.TypeList,
					Description: "An object representing an IPv6 pool.",
					Computed:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"prefix": {
								Type: schema.TypeInt,
								Description: "The number of bits set in the subnet mask." +
									"addresses can be assigned from this pool" +
									"calculated as 2 128-prefix.",
								Computed: true,
							},
							"range": {
								Type:        schema.TypeString,
								Description: "The IPv6 range of addresses in this pool.",
								Computed:    true,
							},
							"region": {
								Type:        schema.TypeString,
								Description: "The region for this pool of IPv6 addresses.",
								Computed:    true,
							},
							"route_target": {
								Type:        schema.TypeString,
								Description: "The last address in this block of IPv6 addresses.",
								Computed:    true,
							},
						},
					},
				},
				"link_local": {
					Type:        schema.TypeList,
					Description: "A link-local IPv6 address that exists in Linode’s system.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
				"slaac": {
					Type:        schema.TypeList,
					Description: "A SLAAC IPv6 address that exists in Linode’s system.",
					Computed:    true,
					Elem:        &schema.Resource{Schema: networkSchema},
				},
			},
		},
	},
}
