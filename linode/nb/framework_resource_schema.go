package nb

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Version: 1,
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the Linode NodeBalancer.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode NodeBalancer.",
			Optional:    true,
		},
		"region": schema.StringAttribute{
			Description:   "The region where this NodeBalancer will be deployed.",
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Default:       stringdefault.StaticString("us-east"),
		},
		"client_conn_throttle": schema.Int64Attribute{
			Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
			Validators: []validator.Int64{
				int64validator.Between(0, 20),
			},
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(0),
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
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
	},
	Blocks: map[string]schema.Block{
		"transfer": schema.ListNestedBlock{
			Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"in": schema.Float64Attribute{
						Description: "The total transfer, in MB, used by this NodeBalancer this month",
						Computed:    true,
					},
					"out": schema.Float64Attribute{
						Description: "The total inbound transfer, in MB, used for this NodeBalancer this month",
						Computed:    true,
					},
					"total": schema.Float64Attribute{
						Description: "The total outbound transfer, in MB, used for this NodeBalancer this month",
						Computed:    true,
					},
				},
			},
		},
	},
}

var resourceNodebalancerV0 = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the Linode NodeBalancer.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the Linode NodeBalancer.",
			Optional:    true,
		},
		"region": schema.StringAttribute{
			Description:   "The region where this NodeBalancer will be deployed.",
			Optional:      true,
			Computed:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			Default:       stringdefault.StaticString("us-east"),
		},
		"client_conn_throttle": schema.Int64Attribute{
			Description: "Throttle connections per second (0-20). Set to 0 (zero) to disable throttling.",
			Validators: []validator.Int64{
				int64validator.Between(0, 20),
			},
			Optional: true,
			Computed: true,
			Default:  int64default.StaticInt64(0),
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
		"tags": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
		},
		"transfer": schema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
	},
}
