package nbvpc

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var DataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"nodebalancer_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer.",
			Required:    true,
		},
		"id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer-VPC configuration.",
			Required:    true,
		},
		"ipv4_range": schema.StringAttribute{
			Description: "A CIDR range for the VPC's IPv4 addresses. " +
				"The NodeBalancer sources IP addresses from this range when " +
				"routing traffic to the backend VPC nodes.",
			Computed: true,
		},
		"ipv6_range": schema.StringAttribute{
			Description: "A CIDR range for the VPC's IPv6 addresses allocated " +
				"as the NodeBalancer's frontend IPs.",
			Computed: true,
		},
		"subnet_id": schema.Int64Attribute{
			Description: "The ID of this configuration's VPC subnet.",
			Computed:    true,
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The ID of this configuration's VPC.",
			Computed:    true,
		},
		"purpose": schema.StringAttribute{
			Description: "Indicates whether the VPC configuration applies to " +
				"backend nodes that serve requests or to the NodeBalancer frontend, " +
				"which manages incoming traffic.",
			Computed: true,
		},
	},
}
