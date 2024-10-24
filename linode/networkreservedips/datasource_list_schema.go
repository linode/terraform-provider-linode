package networkreservedips

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReservedIPObject struct {
	ID         types.String `tfsdk:"id"`
	Address    types.String `tfsdk:"address"`
	Region     types.String `tfsdk:"region"`
	Gateway    types.String `tfsdk:"gateway"`
	SubnetMask types.String `tfsdk:"subnet_mask"`
	Prefix     types.Int64  `tfsdk:"prefix"`
	Type       types.String `tfsdk:"type"`
	Public     types.Bool   `tfsdk:"public"`
	RDNS       types.String `tfsdk:"rdns"`
	LinodeID   types.Int64  `tfsdk:"linode_id"`
	Reserved   types.Bool   `tfsdk:"reserved"`
}

var reservedIPObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"address":     types.StringType,
		"region":      types.StringType,
		"gateway":     types.StringType,
		"subnet_mask": types.StringType,
		"prefix":      types.Int64Type,
		"type":        types.StringType,
		"public":      types.BoolType,
		"rdns":        types.StringType,
		"linode_id":   types.Int64Type,
		"reserved":    types.BoolType,
	},
}

var frameworkDataSourceListSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"reserved_ips": schema.ListAttribute{
			Description: "A list of all reserved IPs.",
			Computed:    true,
			ElementType: reservedIPObjectType,
		},
	},
}
