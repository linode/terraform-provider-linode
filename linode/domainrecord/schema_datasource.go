package domainrecord

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID assigned to this domain record.",
		Optional:    true,
	},
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the Record.",
		Optional:    true,
	},
	"domain_id": {
		Type:        schema.TypeInt,
		Description: "The associated domain's ID.",
		Required:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The type of Record this is in the DNS system.",
		Computed:    true,
	},
	"ttl_sec": {
		Type: schema.TypeInt,
		Description: "The amount of time in seconds that this Domain's records may be cached by resolvers or " +
			"other domain servers.",
		Computed: true,
	},
	"target": {
		Type: schema.TypeString,
		Description: "The target for this Record. This field's actual usage depends on the type of record " +
			"this represents. For A and AAAA records, this is the address the named Domain should resolve to.",
		Computed: true,
	},
	"priority": {
		Type:        schema.TypeInt,
		Description: "The priority of the target host. Lower values are preferred.",
		Computed:    true,
	},
	"weight": {
		Type:        schema.TypeInt,
		Description: "The relative weight of this Record. Higher values are preferred.",
		Computed:    true,
	},
	"port": {
		Type:        schema.TypeInt,
		Description: "The port this Record points to.",
		Computed:    true,
	},
	"protocol": {
		Type:        schema.TypeString,
		Description: "The protocol this Record's service communicates with. Only valid for SRV records.",
		Computed:    true,
	},
	"service": {
		Type:        schema.TypeString,
		Description: "The service this Record identified. Only valid for SRV records.",
		Computed:    true,
	},
	"tag": {
		Type:        schema.TypeString,
		Description: "The tag portion of a CAA record.",
		Computed:    true,
	},
}
