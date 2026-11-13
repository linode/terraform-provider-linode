package natgateway

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceVPCSubnetAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique ID of the attached VPC subnet.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of the attached resource. Always vpc_subnet.",
		Computed:    true,
	},
	"url": schema.StringAttribute{
		Description: "The API URL of the attached VPC subnet.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the VPC subnet.",
		Computed:    true,
	},
	"vpc_id": schema.StringAttribute{
		Description: "The unique ID of the parent VPC.",
		Computed:    true,
	},
	"vpc_label": schema.StringAttribute{
		Description: "The label of the parent VPC.",
		Computed:    true,
	},
}

var dataSourceAddressNestedObject = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"address": schema.StringAttribute{
			Description: "The public IP address assigned to this NAT Gateway.",
			Computed:    true,
		},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The unique ID of the NAT Gateway.",
			Required:    true,
		},
		"region": schema.StringAttribute{
			Description: "The region where the NAT Gateway is located.",
			Computed:    true,
		},
		"label": schema.StringAttribute{
			Description: "The label of the NAT Gateway.",
			Computed:    true,
		},
		"default_ports_per_interface": schema.Int64Attribute{
			Description: "The default number of ports per interface for the NAT Gateway.",
			Computed:    true,
		},
		"addresses": schema.SetNestedAttribute{
			Description:  "The public IP addresses assigned to this NAT Gateway.",
			Computed:     true,
			NestedObject: dataSourceAddressNestedObject,
		},
		"address_autoscale_max": schema.Int64Attribute{
			Description: "The maximum number of auto-allocated IP addresses allowed for this NAT Gateway instance.",
			Computed:    true,
		},
		"portset_assignments": schema.Int64Attribute{
			Description: "The number of portsets currently assigned to interfaces behind this NAT Gateway, excluding system-reserved ports.",
			Computed:    true,
		},
		"portset_capacity": schema.Int64Attribute{
			Description: "The total number of portsets available to allocate across all current and potential autoscaled public IP addresses assigned to this NAT Gateway, excluding system-reserved ports.",
			Computed:    true,
		},
		"vpc_subnet": schema.SingleNestedAttribute{
			Description: "Information about the VPC subnet attached to this NAT Gateway, if there is one attached.",
			Computed:    true,
			Attributes:  dataSourceVPCSubnetAttributes,
		},
		"created": schema.StringAttribute{
			Description: "When this NAT Gateway was created.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
		"updated": schema.StringAttribute{
			Description: "When this NAT Gateway was last updated.",
			Computed:    true,
			CustomType:  timetypes.RFC3339Type{},
		},
	},
}
