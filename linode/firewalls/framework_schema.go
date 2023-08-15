package firewalls

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

var filterConfig = frameworkfilter.Config{
	"id":    {APIFilterable: true, TypeFunc: helper.FilterTypeInt},
	"label": {APIFilterable: true, TypeFunc: helper.FilterTypeString},
	"tags":  {APIFilterable: true, TypeFunc: helper.FilterTypeString},

	"status":  {APIFilterable: false, TypeFunc: helper.FilterTypeString},
	"created": {APIFilterable: false, TypeFunc: helper.FilterTypeString, AllowOrderOverride: true},
	"updated": {APIFilterable: false, TypeFunc: helper.FilterTypeString, AllowOrderOverride: true},
}

var firewallDeviceObject = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the Firewall Device.",
			Computed:    true,
		},
		"entity_id": schema.Int64Attribute{
			Description: "The ID of the underlying entity this device references (i.e. the Linode's ID).",
			Computed:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of Firewall Device.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the underlying entity this device references.",
			Computed:    true,
		},
		"url": schema.StringAttribute{
			Description: "The URL of the underlying entity this device references.",
			Computed:    true,
		},
	},
}

var firewallRuleObject = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label of this rule for display purposes only.",
			Computed:    true,
		},
		"action": schema.StringAttribute{
			Description: "Controls whether traffic is accepted or dropped by this rule (ACCEPT, DROP).",
			Computed:    true,
		},
		"protocol": schema.StringAttribute{
			Description: "The network protocol this rule controls. (TCP, UDP, ICMP)",
			Computed:    true,
		},
		"ports": schema.StringAttribute{
			Description: "A string representation of ports and/or port ranges (i.e. \"443\" or \"80-90, 91\").",
			Computed:    true,
		},
		"ipv4": schema.SetAttribute{
			Description: "A list of IPv4 addresses or networks in IP/mask format.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"ipv6": schema.SetAttribute{
			Description: "A list of IPv6 addresses or networks in IP/mask format.",
			ElementType: types.StringType,
			Computed:    true,
		},
	},
}

var firewallObject = schema.NestedBlockObject{
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
		"disabled": schema.BoolAttribute{
			Description: "If true, the Firewall is inactive.",
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
		"linodes": schema.SetAttribute{
			ElementType: types.Int64Type,
			Description: "The IDs of Linodes to apply this firewall to.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The status of the firewall.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			CustomType:  customtypes.RFC3339TimeStringType{},
			Description: "When this Firewall was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			CustomType:  customtypes.RFC3339TimeStringType{},
			Description: "When this Firewall was last updated.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"inbound": schema.ListNestedBlock{
			Description:  "A set of firewall rules that specify what inbound network traffic is allowed.",
			NestedObject: firewallRuleObject,
		},
		"outbound": schema.ListNestedBlock{
			Description:  "A set of firewall rules that specify what outbound network traffic is allowed.",
			NestedObject: firewallRuleObject,
		},
		"devices": schema.ListNestedBlock{
			Description:  "The devices associated with this firewall.",
			NestedObject: firewallDeviceObject,
		},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"order":    filterConfig.OrderSchema(),
		"order_by": filterConfig.OrderBySchema(),
	},
	Blocks: map[string]schema.Block{
		"filter": filterConfig.Schema(),
		"firewalls": schema.ListNestedBlock{
			Description:  "The returned list of Firewalls.",
			NestedObject: firewallObject,
		},
	},
}
