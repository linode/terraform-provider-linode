package reservedips

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
)

type ReservedIPFilterModel struct {
	ID          types.String                     `tfsdk:"id"`
	Filters     frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order       types.String                     `tfsdk:"order"`
	OrderBy     types.String                     `tfsdk:"order_by"`
	ReservedIPs []ReservedIPObject               `tfsdk:"reserved_ips"`
}

func (data *ReservedIPFilterModel) parseReservedIPs(
	ctx context.Context,
	ips []linodego.InstanceIP,
) diag.Diagnostics {
	result := make([]ReservedIPObject, len(ips))
	for i, ip := range ips {
		var ipData ReservedIPObject
		ipData.ID = types.StringValue(ip.Address)
		ipData.Address = types.StringValue(ip.Address)
		ipData.Region = types.StringValue(ip.Region)
		ipData.Gateway = types.StringValue(ip.Gateway)
		ipData.SubnetMask = types.StringValue(ip.SubnetMask)
		ipData.Prefix = types.Int64Value(int64(ip.Prefix))
		ipData.Type = types.StringValue(string(ip.Type))
		ipData.Public = types.BoolValue(ip.Public)
		ipData.RDNS = types.StringValue(ip.RDNS)
		ipData.LinodeID = types.Int64Value(int64(ip.LinodeID))
		ipData.Reserved = types.BoolValue(ip.Reserved)

		vpcNAT1To1List, diags := types.ListValueFrom(ctx, instancenetworking.VPCNAT1To1Type, ip.VPCNAT1To1)
		if diags.HasError() {
			return diags
		}
		ipData.IPVPCNAT1To1 = vpcNAT1To1List

		result[i] = ipData
	}

	data.ReservedIPs = result

	return nil
}
