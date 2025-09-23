package networkingips

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
)

type IPAddressModel struct {
	Address     types.String `tfsdk:"address"`
	Type        types.String `tfsdk:"type"`
	Region      types.String `tfsdk:"region"`
	RDNS        types.String `tfsdk:"rdns"`
	Prefix      types.Int64  `tfsdk:"prefix"`
	Gateway     types.String `tfsdk:"gateway"`
	SubnetMask  types.String `tfsdk:"subnet_mask"`
	Public      types.Bool   `tfsdk:"public"`
	LinodeID    types.Int64  `tfsdk:"linode_id"`
	InterfaceID types.Int64  `tfsdk:"interface_id"`
	Reserved    types.Bool   `tfsdk:"reserved"`
	VPCNAT1To1  types.Object `tfsdk:"vpc_nat_1_1"`
}

func (m *IPAddressModel) ParseIP(ip linodego.InstanceIP) diag.Diagnostics {
	m.Address = types.StringValue(ip.Address)
	m.Type = types.StringValue(string(ip.Type))
	m.Region = types.StringValue(ip.Region)
	m.RDNS = types.StringValue(ip.RDNS)
	m.Prefix = types.Int64Value(int64(ip.Prefix))
	m.Gateway = types.StringValue(ip.Gateway)
	m.SubnetMask = types.StringValue(ip.SubnetMask)
	m.Public = types.BoolValue(ip.Public)
	m.LinodeID = types.Int64Value(int64(ip.LinodeID))
	m.InterfaceID = types.Int64PointerValue(helper.IntPtrToInt64Ptr(ip.InterfaceID))
	m.Reserved = types.BoolValue(ip.Reserved)

	vpcNAT1To1, d := instancenetworking.FlattenIPVPCNAT1To1(ip.VPCNAT1To1)
	if d.HasError() {
		return d
	}

	m.VPCNAT1To1 = vpcNAT1To1

	return nil
}

// FilterModel describes the Terraform resource data model to match the
// resource schema.
type FilterModel struct {
	ID          types.String                     `tfsdk:"id"`
	Filters     frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order       types.String                     `tfsdk:"order"`
	OrderBy     types.String                     `tfsdk:"order_by"`
	IPAddresses []IPAddressModel                 `tfsdk:"ip_addresses"`
}

func (data *FilterModel) parseIPAddresses(
	ips []linodego.InstanceIP,
) (d diag.Diagnostics) {
	result := make([]IPAddressModel, len(ips))

	for i := range ips {
		var data IPAddressModel
		d.Append(data.ParseIP(ips[i])...)
		result[i] = data
	}

	data.IPAddresses = result

	return d
}
