package linodeinterface

import (
	"github.com/hashicorp/terraform-plugin-framework-nettypes/cidrtypes"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	dataSourceDefaultRouteAttribute = schema.SingleNestedAttribute{
		Description: "Default route configuration for the interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"ipv4": schema.BoolAttribute{
				Description: "Whether this interface is used for the IPv4 default route.",
				Computed:    true,
			},
			"ipv6": schema.BoolAttribute{
				Description: "Whether this interface is used for the IPv6 default route.",
				Computed:    true,
			},
		},
	}

	dataSourcePublicIPv4AddressAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Description: "The IPv4 address.",
				Computed:    true,
			},
			"primary": schema.BoolAttribute{
				Description: "Whether this is the primary IPv4 address.",
				Computed:    true,
			},
		},
	}

	dataSourceSharedPublicIPv4AddressAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Description: "The shared IPv4 address.",
				Computed:    true,
			},
			"linode_id": schema.Int64Attribute{
				Description: "The ID of the Linode that this shared address belongs to.",
				Computed:    true,
			},
		},
	}

	dataSourcePublicIPv6RangeAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"range": schema.StringAttribute{
				CustomType:  cidrtypes.IPv6PrefixType{},
				Description: "The IPv6 range.",
				Computed:    true,
			},
			"route_target": schema.StringAttribute{
				Description: "The route target for this IPv6 range.",
				Computed:    true,
			},
		},
	}

	dataSourcePublicIPv6SLAACAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Description: "The IPv6 SLAAC address.",
				CustomType:  iptypes.IPv6AddressType{},
				Computed:    true,
			},
			"prefix": schema.Int64Attribute{
				Description: "The prefix length for the IPv6 SLAAC address.",
				Computed:    true,
			},
		},
	}

	dataSourcePublicIPv4Attribute = schema.SingleNestedAttribute{
		Description: "The public IPv4 configuration for the interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"addresses": schema.SetNestedAttribute{
				Description:  "IPv4 addresses assigned to this interface.",
				Computed:     true,
				NestedObject: dataSourcePublicIPv4AddressAttribute,
			},
			"shared": schema.SetNestedAttribute{
				Description:  "IPv4 addresses shared with other Linodes.",
				Computed:     true,
				NestedObject: dataSourceSharedPublicIPv4AddressAttribute,
			},
		},
	}

	dataSourcePublicIPv6Attribute = schema.SingleNestedAttribute{
		Description: "The public IPv6 configuration for the interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"ranges": schema.SetNestedAttribute{
				Description:  "IPv6 ranges assigned to this interface.",
				Computed:     true,
				NestedObject: dataSourcePublicIPv6RangeAttribute,
			},
			"shared": schema.SetNestedAttribute{
				Description:  "IPv6 ranges shared with other Linodes.",
				Computed:     true,
				NestedObject: dataSourcePublicIPv6RangeAttribute,
			},
			"slaac": schema.SetNestedAttribute{
				Description:  "IPv6 SLAAC configuration.",
				Computed:     true,
				NestedObject: dataSourcePublicIPv6SLAACAttribute,
			},
		},
	}

	dataSourcePublicAttribute = schema.SingleNestedAttribute{
		Description: "Configuration profile for the public interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"ipv4": dataSourcePublicIPv4Attribute,
			"ipv6": dataSourcePublicIPv6Attribute,
		},
	}

	dataSourceVLANAttribute = schema.SingleNestedAttribute{
		Description: "Configuration profile for the VLAN interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"vlan_label": schema.StringAttribute{
				Description: "The label of the VLAN.",
				Computed:    true,
			},
			"ipam_address": schema.StringAttribute{
				CustomType:  cidrtypes.IPv4PrefixType{},
				Description: "The IPAM (IP Address Management) address of the VLAN interface.",
				Computed:    true,
			},
		},
	}

	dataSourceVPCIPv4AddressAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Description: "The VPC IPv4 address.",
				CustomType:  iptypes.IPv4AddressType{},
				Computed:    true,
			},
			"primary": schema.BoolAttribute{
				Description: "Whether this is the primary VPC IPv4 address.",
				Computed:    true,
			},
			"nat_1_1_address": schema.StringAttribute{
				Description: "The 1:1 NAT address for this VPC IPv4 address.",
				CustomType:  iptypes.IPv4AddressType{},
				Computed:    true,
			},
		},
	}

	dataSourceVPCIPv4RangeAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"range": schema.StringAttribute{
				CustomType:  cidrtypes.IPv4PrefixType{},
				Description: "The VPC IPv4 range.",
				Computed:    true,
			},
		},
	}

	dataSourceVPCIPv6SLAACAttribute = schema.NestedAttributeObject{
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

	dataSourceVPCIPv6RangeAttribute = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"range": schema.StringAttribute{
				Description: "The IPv6 network range in CIDR notation.",
				CustomType:  cidrtypes.IPv6PrefixType{},
				Computed:    true,
			},
		},
	}

	dataSourceVPCIPv4Attribute = schema.SingleNestedAttribute{
		Description: "The IPv4 configuration for the VPC interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"addresses": schema.SetNestedAttribute{
				Description:  "IPv4 addresses assigned to this VPC interface.",
				Computed:     true,
				NestedObject: dataSourceVPCIPv4AddressAttribute,
			},
			"ranges": schema.SetNestedAttribute{
				Description:  "IPv4 ranges assigned to this VPC interface.",
				Computed:     true,
				NestedObject: dataSourceVPCIPv4RangeAttribute,
			},
		},
	}

	dataSourceVPCIPv6Attribute = schema.SingleNestedAttribute{
		Description: "The IPv6 configuration for the VPC interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"is_public": schema.BoolAttribute{
				Description: "Indicates whether the IPv6 configuration on the Linode interface is public.",
				Computed:    true,
			},
			"slaac": schema.SetNestedAttribute{
				Description:  "IPv6 SLAAC address ranges.",
				Computed:     true,
				NestedObject: dataSourceVPCIPv6SLAACAttribute,
			},
			"ranges": schema.SetNestedAttribute{
				Description:  "IPv6 ranges assigned to this VPC interface.",
				Computed:     true,
				NestedObject: dataSourceVPCIPv6RangeAttribute,
			},
		},
	}

	dataSourceVPCAttribute = schema.SingleNestedAttribute{
		Description: "Configuration profile for the VPC interface.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"subnet_id": schema.Int64Attribute{
				Description: "The ID of the VPC subnet.",
				Computed:    true,
			},
			"ipv4": dataSourceVPCIPv4Attribute,
			"ipv6": dataSourceVPCIPv6Attribute,
		},
	}
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Linode interface.",
			Required:    true,
		},
		"linode_id": schema.Int64Attribute{
			Description: "The ID of the Linode.",
			Required:    true,
		},
		"default_route": dataSourceDefaultRouteAttribute,
		"public":        dataSourcePublicAttribute,
		"vlan":          dataSourceVLANAttribute,
		"vpc":           dataSourceVPCAttribute,
	},
}
