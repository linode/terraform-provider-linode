package firewall

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var deviceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":        types.Int64Type,
		"entity_id": types.Int64Type,
		"type":      types.StringType,
		"label":     types.StringType,
		"url":       types.StringType,
	},
}

var RuleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"label":    types.StringType,
		"action":   types.StringType,
		"ports":    types.StringType,
		"protocol": types.StringType,
		"ipv4":     types.ListType{ElemType: types.StringType},
		"ipv6":     types.ListType{ElemType: types.StringType},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID assigned to this Firewall.",
			Required:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label for the Firewall. For display purposes only. If no label is provided, a " +
				"default will be assigned.",
			Computed: true,
		},
		"tags": schema.SetAttribute{
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"disabled": schema.BoolAttribute{
			Description: "If true, the Firewall is inactive.",
			Computed:    true,
		},
		"inbound": schema.ListAttribute{
			ElementType: RuleObjectType,
			Description: "A firewall rule that specifies what inbound network traffic is allowed.",
			Computed:    true,
		},
		"inbound_policy": schema.StringAttribute{
			Description: "The default behavior for inbound traffic. This setting can be overridden by updating " +
				"the inbound.action property for an individual Firewall Rule.",
			Computed: true,
		},
		"outbound": schema.ListAttribute{
			ElementType: RuleObjectType,
			Description: "A firewall rule that specifies what outbound network traffic is allowed.",
			Computed:    true,
		},
		"outbound_policy": schema.StringAttribute{
			Description: "The default behavior for outbound traffic. This setting can be overridden by updating " +
				"the outbound.action property for an individual Firewall Rule.",
			Computed: true,
		},
		"linodes": schema.SetAttribute{
			ElementType: types.Int64Type,
			Description: "The IDs of Linodes assigned to this Firewall.",
			Computed:    true,
		},
		"nodebalancers": schema.SetAttribute{
			ElementType: types.Int64Type,
			Description: "The IDs of NodeBalancers assigned to this Firewall.",
			Computed:    true,
		},
		"devices": schema.ListAttribute{
			ElementType: deviceObjectType,
			Description: "The devices associated with this firewall.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the firewall.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "When this Firewall was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "When this Firewall was last updated.",
			Computed:    true,
		},
	},
}
