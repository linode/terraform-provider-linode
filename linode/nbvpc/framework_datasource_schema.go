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
			Computed:    true,
		},

		"ipv4_range": schema.StringAttribute{
			Description: "A CIDR range for the VPC's IPv4 addresses. " +
				"The NodeBalancer sources IP addresses from this range when " +
				"routing traffic to the backend VPC nodes.",
			Computed: true,
		},
		"ipv6_range": schema.StringAttribute{
			Description: "A CIDR range for the VPC's IPv6 addresses. " +
				"The NodeBalancer sources IP addresses from this " +
				"range when routing traffic to the backend VPC nodes.",
			Computed: true,
		},

		"subnet_id": schema.Int64Attribute{
			Description: "The VPC's subnet. Run the List VPCs operation to retrieve VPCs and their subnets.",
			Computed:    true,
		},
		"vpc_id": schema.Int64Attribute{
			Description: "The `id` of the VPC configured for this NodeBalancer.",
			Computed:    true,
		},
	},
}
