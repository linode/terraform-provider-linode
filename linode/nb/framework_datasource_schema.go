package nb

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var transferObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"in":    types.Float64Type,
		"out":   types.Float64Type,
		"total": types.Float64Type,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the Linode NodeBalancer.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode NodeBalancer.",
			Computed:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region where this NodeBalancer will be deployed.",
			Computed:    true,
		},
		"client_conn_throttle": schema.Int64Attribute{
			Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
			Computed:    true,
		},
		"hostname": schema.StringAttribute{
			Description: "This NodeBalancer's hostname, ending with .nodebalancer.linode.com",
			Computed:    true,
		},
		"ipv4": schema.StringAttribute{
			Description: "The Public IPv4 Address of this NodeBalancer",
			Computed:    true,
		},
		"ipv6": schema.StringAttribute{
			Description: "The Public IPv6 Address of this NodeBalancer",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this NodeBalancer was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this NodeBalancer was last updated.",
			Computed:    true,
		},
		"transfer": schema.ListAttribute{
			Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
			Computed:    true,
			ElementType: transferObjectType,
		},
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
	},
}
