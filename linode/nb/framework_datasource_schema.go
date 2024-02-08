package nb

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalls"
)

var TransferObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"in":    types.Float64Type,
		"out":   types.Float64Type,
		"total": types.Float64Type,
	},
}

var NodeBalancerAttributes = map[string]schema.Attribute{
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
		CustomType:  timetypes.RFC3339Type{},
	},
	"updated": schema.StringAttribute{
		Description: "When this NodeBalancer was last updated.",
		Computed:    true,
		CustomType:  timetypes.RFC3339Type{},
	},
	"transfer": schema.ListAttribute{
		Description: "Information about the amount of transfer this NodeBalancer has had so far this month.",
		Computed:    true,
		ElementType: TransferObjectType,
	},
	"tags": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: NodeBalancerAttributes,
	Blocks: map[string]schema.Block{
		"firewalls": schema.ListNestedBlock{
			Description: "A list of Firewalls assigned to this NodeBalancer.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The unique ID assigned to this Firewall.",
						Computed:    true,
					},
					"label": schema.StringAttribute{
						Computed: true,
						Description: "The label for the Firewall. For display purposes only. If no label is provided, a " +
							"default will be assigned.",
					},
					"tags": schema.SetAttribute{
						Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"inbound_policy": schema.StringAttribute{
						Description: "The default behavior for inbound traffic.",
						Computed:    true,
					},
					"outbound_policy": schema.StringAttribute{
						Description: "The default behavior for outbound traffic.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of the firewall.",
						Computed:    true,
					},
					"created": schema.StringAttribute{
						CustomType:  timetypes.RFC3339Type{},
						Description: "When this Firewall was created.",
						Computed:    true,
					},
					"updated": schema.StringAttribute{
						CustomType:  timetypes.RFC3339Type{},
						Description: "When this Firewall was last updated.",
						Computed:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"inbound": schema.ListNestedBlock{
						Description:  "A set of firewall rules that specify what inbound network traffic is allowed.",
						NestedObject: firewalls.FirewallRuleObject,
					},
					"outbound": schema.ListNestedBlock{
						Description:  "A set of firewall rules that specify what outbound network traffic is allowed.",
						NestedObject: firewalls.FirewallRuleObject,
					},
				},
			},
		},
	},
}
