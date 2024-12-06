package reservedips

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

type ReservedIPObject struct {
	ID           types.String `tfsdk:"id"`
	Address      types.String `tfsdk:"address"`
	Region       types.String `tfsdk:"region"`
	Gateway      types.String `tfsdk:"gateway"`
	SubnetMask   types.String `tfsdk:"subnet_mask"`
	Prefix       types.Int64  `tfsdk:"prefix"`
	Type         types.String `tfsdk:"type"`
	Public       types.Bool   `tfsdk:"public"`
	RDNS         types.String `tfsdk:"rdns"`
	LinodeID     types.Int64  `tfsdk:"linode_id"`
	Reserved     types.Bool   `tfsdk:"reserved"`
	IPVPCNAT1To1 types.List   `tfsdk:"vpc_nat_1_1"`
}

var ReservedIPAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "The unique identifier of the reserved IP.",
		Computed:    true,
	},
	"address": schema.StringAttribute{
		Description: "The IP address.",
		Computed:    true,
	},
	"region": schema.StringAttribute{
		Description: "The region where the IP is located.",
		Optional:    true,
		Computed:    true,
	},
	"gateway": schema.StringAttribute{
		Description: "The gateway for the reserved IP.",
		Computed:    true,
	},
	"subnet_mask": schema.StringAttribute{
		Description: "The subnet mask for the reserved IP.",
		Computed:    true,
	},
	"prefix": schema.Int64Attribute{
		Description: "The prefix length of the reserved IP.",
		Computed:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of the reserved IP.",
		Computed:    true,
	},
	"public": schema.BoolAttribute{
		Description: "Whether the IP is public.",
		Computed:    true,
	},
	"rdns": schema.StringAttribute{
		Description: "The reverse DNS for the reserved IP.",
		Optional:    true,
		Computed:    true,
	},
	"linode_id": schema.Int64Attribute{
		Description: "The Linode ID associated with this reserved IP.",
		Computed:    true,
	},
	"reserved": schema.BoolAttribute{
		Description: "Indicates if this IP is reserved.",
		Computed:    true,
	},
	"vpc_nat_1_1": schema.ListAttribute{
		ElementType: instancenetworking.VPCNAT1To1Type,
		Computed:    true,
	},
}

var filterConfig = frameworkfilter.Config{
	"region": {APIFilterable: true, TypeFunc: frameworkfilter.FilterTypeString},
}

var frameworkDataSourceSchema = schema.Schema{
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
		"reserved_ips": schema.ListNestedBlock{
			Description: "The returned list of Reserved IPs.",
			NestedObject: schema.NestedBlockObject{
				Attributes: ReservedIPAttributes,
			},
		},
	},
}
