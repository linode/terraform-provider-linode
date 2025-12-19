package linodeinterface

import (
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	linodesetplanmodifier "github.com/linode/terraform-provider-linode/v3/linode/helper/setplanmodifiers"
)

var configuredPublicInterfaceIPv4Address = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("auto"),
		},
		"primary": schema.BoolAttribute{
			Optional: true,
		},
	},
}

var computedPublicInterfaceIPv4Address = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			CustomType: iptypes.IPv4AddressType{},
			Computed:   true,
		},
		"primary": schema.BoolAttribute{
			Computed: true,
		},
	},
}

var sharedPublicInterfaceIPv4Address = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			CustomType: iptypes.IPv4AddressType{},
			Computed:   true,
		},
		"linode_id": schema.Int64Attribute{
			Computed: true,
		},
	},
}

var configuredPublicInterfaceIPv6Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Required: true,
		},
	},
}

var computedPublicInterfaceIPv6Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			CustomType: cidrtypes.IPv6PrefixType{},
			Computed:   true,
		},
		"route_target": schema.StringAttribute{
			Description: "The public IPv6 address that the range is routed to.",
			CustomType:  iptypes.IPv6AddressType{},
			Computed:    true,
		},
	},
}

var resourcePublicInterfaceIPv6SLAAC = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Computed:   true,
			CustomType: iptypes.IPv6AddressType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"prefix": schema.Int64Attribute{
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
	},
}

var configuredVPCInterfaceIPv4Address = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString("auto"),
		},
		"primary": schema.BoolAttribute{
			Optional: true,
		},
		"nat_1_1_address": schema.StringAttribute{
			Description: "The 1:1 NAT IPv4 address used to associate a public " +
				"IPv4 address with the interface's VPC subnet IPv4 address.",
			Optional: true,
		},
	},
}

var computedVPCInterfaceIPv4Address = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Computed: true,
		},
		"primary": schema.BoolAttribute{
			Computed: true,
		},
		"nat_1_1_address": schema.StringAttribute{
			Description: "The assigned 1:1 NAT IPv4 address used to associate " +
				"a public IPv4 address with the interface's VPC subnet IPv4 " +
				"address, calculated from `nat_1_1_address`.",
			Computed: true,
		},
	},
}

var configuredVPCInterfaceIPv4Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Required: true,
		},
	},
}

var computedVPCInterfaceIPv4Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			CustomType: cidrtypes.IPv4PrefixType{},
			Computed:   true,
		},
	},
}

var configuredVPCInterfaceIPv6SLAAC = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 network range in CIDR notation.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("auto"),
		},
	},
}

var computedVPCInterfaceIPv6SLAAC = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 network range in CIDR notation.",
			CustomType:  cidrtypes.IPv6PrefixType{},
			Computed:    true,
		},
		"address": schema.StringAttribute{
			Description: "The assigned IPv6 address within the range.",
			CustomType:  iptypes.IPv6AddressType{},
			Computed:    true,
		},
	},
}

var configuredVPCInterfaceIPv6Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 network range in CIDR notation.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("auto"),
		},
	},
}

var computedVPCInterfaceIPv6Range = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"range": schema.StringAttribute{
			Description: "The IPv6 network range in CIDR notation.",
			CustomType:  cidrtypes.IPv6PrefixType{},
			Computed:    true,
		},
	},
}

var resourcePublicIPv4Attribute = schema.SingleNestedAttribute{
	Description: "IPv4 addresses for this interface.",
	Optional:    true,
	Computed:    true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]schema.Attribute{
		"addresses": schema.ListNestedAttribute{
			Description:  "IPv4 addresses configured for this Linode interface.",
			Optional:     true,
			NestedObject: configuredPublicInterfaceIPv4Address,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_addresses": schema.SetNestedAttribute{
			Description:  "The IPv4 address exclusively assigned to this Linode interface.",
			Computed:     true,
			NestedObject: computedPublicInterfaceIPv4Address,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("public").AtName("ipv4").AtName("addresses"),
				),
			},
		},
		"shared": schema.SetNestedAttribute{
			Description:  "The IPv4 address assigned to this Linode interface, which is also shared with another Linode.",
			Computed:     true,
			NestedObject: sharedPublicInterfaceIPv4Address,
		},
	},
}

var resourcePublicIPv6Attribute = schema.SingleNestedAttribute{
	Description: "IPv6 addresses for this interface.",
	Optional:    true,
	Computed:    true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]schema.Attribute{
		"ranges": schema.ListNestedAttribute{
			Description:  "Configured IPv6 range in CIDR notation (2600:0db8::1/64) or prefix-only (/64).",
			Optional:     true,
			NestedObject: configuredPublicInterfaceIPv6Range,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_ranges": schema.SetNestedAttribute{
			Description:  "The IPv6 ranges exclusively assigned to this Linode interface.",
			Computed:     true,
			NestedObject: computedPublicInterfaceIPv6Range,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("public").AtName("ipv6").AtName("ranges"),
				),
			},
		},
		"shared": schema.SetNestedAttribute{
			Description:  "The IPv6 address assigned to this Linode interface, which is also shared with another Linode.",
			Computed:     true,
			NestedObject: computedPublicInterfaceIPv6Range,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
		},
		"slaac": schema.SetNestedAttribute{
			Description: "The public slaac and subnet prefix settings for this public interface that is used to " +
				"communicate over the public internet, and with other services in the same data center.",
			Computed: true,
			PlanModifiers: []planmodifier.Set{
				setplanmodifier.UseStateForUnknown(),
			},
			NestedObject: resourcePublicInterfaceIPv6SLAAC,
		},
	},
}

var resourcePublicInterfaceAttribute = schema.SingleNestedAttribute{
	Description: "Linode public interface.",
	Optional:    true,
	Validators: []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("vlan"),
			path.MatchRelative().AtParent().AtName("vpc"),
		),
	},
	Attributes: map[string]schema.Attribute{
		"ipv4": resourcePublicIPv4Attribute,
		"ipv6": resourcePublicIPv6Attribute,
	},
}

var resourceVPCIPv4Attribute = schema.SingleNestedAttribute{
	Optional: true,
	Computed: true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]schema.Attribute{
		"addresses": schema.ListNestedAttribute{
			Description:  "Specifies the IPv4 addresses to use in the VPC subnet.",
			Optional:     true,
			NestedObject: configuredVPCInterfaceIPv4Address,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_addresses": schema.SetNestedAttribute{
			Description:  "Assigned IPv4 addresses to use in the VPC subnet, calculated from `addresses` input.",
			Computed:     true,
			NestedObject: computedVPCInterfaceIPv4Address,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("vpc").AtName("ipv4").AtName("addresses"),
				),
			},
		},
		"ranges": schema.ListNestedAttribute{
			Description:  "CIDR notation of a range (1.2.3.4/24) or prefix only (/24).",
			Optional:     true,
			NestedObject: configuredVPCInterfaceIPv4Range,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_ranges": schema.SetNestedAttribute{
			Description:  "Assigned IPv4 ranges to use in the VPC subnet, calculated from `ranges` input.",
			Computed:     true,
			NestedObject: computedVPCInterfaceIPv4Range,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("vpc").AtName("ipv4").AtName("ranges"),
				),
			},
		},
	},
}

var resourceVPCIPv6Attribute = schema.SingleNestedAttribute{
	Optional: true,
	Computed: true,
	PlanModifiers: []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	},
	Attributes: map[string]schema.Attribute{
		"is_public": schema.BoolAttribute{
			Description: "Indicates whether the IPv6 configuration on the Linode interface is public.",
			Optional:    true,
			Computed:    true,
		},
		"slaac": schema.ListNestedAttribute{
			Description:  "Defines IPv6 SLAAC address ranges.",
			Optional:     true,
			NestedObject: configuredVPCInterfaceIPv6SLAAC,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_slaac": schema.SetNestedAttribute{
			Description:  "Assigned IPv6 SLAAC address ranges, calculated from `addresses` input.",
			Computed:     true,
			NestedObject: computedVPCInterfaceIPv6SLAAC,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("vpc").AtName("ipv6").AtName("slaac"),
				),
			},
		},
		"ranges": schema.ListNestedAttribute{
			Description:  "CIDR notation of a range (1.2.3.4/24) or prefix only (/24).",
			Optional:     true,
			NestedObject: configuredVPCInterfaceIPv6Range,
			Validators: []validator.List{
				listvalidator.NoNullValues(),
			},
		},
		"assigned_ranges": schema.SetNestedAttribute{
			Description:  "Assigned IPv6 ranges to use in the VPC subnet, calculated from `ranges` input.",
			Computed:     true,
			NestedObject: computedVPCInterfaceIPv6Range,
			PlanModifiers: []planmodifier.Set{
				linodesetplanmodifier.UseStateForUnknownUnlessTheseChanged(
					path.MatchRoot("vpc").AtName("ipv6").AtName("ranges"),
				),
			},
		},
	},
}

var vpcInterfaceSchema = schema.SingleNestedAttribute{
	Description: "Linode VPC interface.",
	Optional:    true,
	Validators: []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("public"),
			path.MatchRelative().AtParent().AtName("vlan"),
		),
	},
	Attributes: map[string]schema.Attribute{
		"ipv4": resourceVPCIPv4Attribute,
		"ipv6": resourceVPCIPv6Attribute,
		"subnet_id": schema.Int64Attribute{
			Required:    true,
			Description: "The VPC subnet identifier for this interface.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
	},
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID for this interface.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode to assign this interface to.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"firewall_id": schema.Int64Attribute{
			Description: "ID of an enabled firewall to secure a VPC or public interface. Not allowed for VLAN interfaces.",
			Optional:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"default_route": schema.SingleNestedAttribute{
			Description: "Indicates if the interface serves as the default route when multiple interfaces are eligible for this role.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"ipv4": schema.BoolAttribute{
					Description: "If set to true, the interface is used for the IPv4 default route.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"ipv6": schema.BoolAttribute{
					Description: "If set to true, the interface is used for the IPv6 default route.",
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"public": resourcePublicInterfaceAttribute,
		"vlan": schema.SingleNestedAttribute{
			Description: "Linode VLAN interface.",
			Optional:    true,
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRelative().AtParent().AtName("public"),
					path.MatchRelative().AtParent().AtName("vpc"),
				),
			},
			Attributes: map[string]schema.Attribute{
				"ipam_address": schema.StringAttribute{
					Description: "This VLAN interface's private IPv4 address in classless inter-domain routing (CIDR) notation.",
					Optional:    true,
					CustomType:  cidrtypes.IPv4PrefixType{},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"vlan_label": schema.StringAttribute{
					Description: "The VLAN's unique label.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 64),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
		"vpc": vpcInterfaceSchema,
	},
}
