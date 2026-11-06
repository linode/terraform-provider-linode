package natgatewayinterfaces

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var interfaceLinodeAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The ID of the Linode the interface belongs to.",
		Computed:    true,
	},
	"label": schema.StringAttribute{
		Description: "The label of the Linode the interface belongs to.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of the attached resource. Always linode.",
		Computed:    true,
	},
	"url": schema.StringAttribute{
		Description: "The API URL of the Linode the interface belongs to.",
		Computed:    true,
	},
}

var interfacePortsetPortAttributes = map[string]schema.Attribute{
	"start": schema.Int64Attribute{
		Description: "The first port in the range.",
		Computed:    true,
	},
	"end": schema.Int64Attribute{
		Description: "The last port in the range.",
		Computed:    true,
	},
}

var interfacePortsetAttributes = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Description: "The public IPv4 address the portset belongs to.",
		Computed:    true,
	},
	"ports": schema.ListNestedAttribute{
		Description: "The port ranges in this portset.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: interfacePortsetPortAttributes,
		},
	},
}

var interfaceAttributes = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		Description: "The unique ID of the NAT Gateway interface.",
		Computed:    true,
	},
	"linode": schema.SingleNestedAttribute{
		Description: "The Linode the interface belongs to.",
		Computed:    true,
		Attributes:  interfaceLinodeAttributes,
	},
	"addresses": schema.ListAttribute{
		Description: "The public IPv4 addresses assigned to the interface.",
		Computed:    true,
		ElementType: types.StringType,
	},
	"portsets": schema.ListNestedAttribute{
		Description: "The portsets assigned to the interface. May not be available for all accounts.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: interfacePortsetAttributes,
		},
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"natgateway_id": schema.Int64Attribute{
			Description: "The ID of the NAT Gateway to list interfaces for.",
			Required:    true,
		},
		"address": schema.StringAttribute{
			Description: "A public IPv4 address of the NAT Gateway. When set, only interfaces using this address are returned.",
			Optional:    true,
		},
		"interfaces": schema.ListNestedAttribute{
			Description: "The returned list of NAT Gateway interfaces.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: interfaceAttributes,
			},
		},
	},
}
