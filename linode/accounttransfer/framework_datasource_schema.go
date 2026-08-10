package accounttransfer

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var regionTransferAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The Region ID for this network utilization data.",
		Computed:    true,
	},
	"billable": schema.Int64Attribute{
		Description: "The amount of your transfer pool that is billable this billing cycle for this Region.",
		Computed:    true,
	},
	"quota": schema.Int64Attribute{
		Description: "The amount of network usage allowed this billing cycle for this Region.",
		Computed:    true,
	},
	"used": schema.Int64Attribute{
		Description: "The amount of network usage you have used this billing cycle for this Region.",
		Computed:    true,
	},
}

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The data source's unique ID.",
			Computed:    true,
		},
		"billable": schema.Int64Attribute{
			Description: "The amount of your transfer pool that is billable this billing cycle.",
			Computed:    true,
		},
		"quota": schema.Int64Attribute{
			Description: "The amount of network usage allowed this billing cycle.",
			Computed:    true,
		},
		"used": schema.Int64Attribute{
			Description: "The amount of network usage you have used this billing cycle.",
			Computed:    true,
		},
		"region_transfers": schema.ListNestedAttribute{
			Description: "The network utilization for the current month in regions with separate utilization quotas and rates.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: regionTransferAttributes,
			},
		},
	},
}
