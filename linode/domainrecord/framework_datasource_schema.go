package domainrecord

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID assigned to this domain record.",
			Optional:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the Record.",
			Optional:    true,
		},
		"domain_id": schema.Int64Attribute{
			Description: "The associated domain's ID.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of Record this is in the DNS system.",
			Computed:    true,
		},
		"ttl_sec": schema.Int64Attribute{
			Description: "The amount of time in seconds that this Domain's records may be cached by resolvers or " +
				"other domain servers.",
			Computed: true,
		},
		"target": schema.StringAttribute{
			Description: "The target for this Record. This field's actual usage depends on the type of record " +
				"this represents. For A and AAAA records, this is the address the named Domain should resolve to.",
			Computed: true,
		},
		"priority": schema.Int64Attribute{
			Description: "The priority of the target host. Lower values are preferred.",
			Computed:    true,
		},
		"weight": schema.Int64Attribute{
			Description: "The relative weight of this Record. Higher values are preferred.",
			Computed:    true,
		},
		"port": schema.Int64Attribute{
			Description: "The port this Record points to.",
			Computed:    true,
		},
		"protocol": schema.StringAttribute{
			Description: "The protocol this Record's service communicates with. Only valid for SRV records.",
			Computed:    true,
		},
		"service": schema.StringAttribute{
			Description: "The service this Record identified. Only valid for SRV records.",
			Computed:    true,
		},
		"tag": schema.StringAttribute{
			Description: "The tag portion of a CAA record.",
			Computed:    true,
		},
	},
}
