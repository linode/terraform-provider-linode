package vpcsubnets

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

// VPCSubnetFilterModel describes the Terraform resource data model to match the
// resource schema.
type VPCSubnetFilterModel struct {
	ID         types.String                     `tfsdk:"id"`
	VPCId      types.Int64                      `tfsdk:"vpc_id"`
	Filters    frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order      types.String                     `tfsdk:"order"`
	OrderBy    types.String                     `tfsdk:"order_by"`
	VPCSubnets []VPCSubnetModel                 `tfsdk:"vpc_subnets"`
}

type VPCSubnetModel struct {
	ID      types.Int64                        `tfsdk:"id"`
	Label   types.String                       `tfsdk:"label"`
	IPv4    types.String                       `tfsdk:"ipv4"`
	Linodes types.List                         `tfsdk:"linodes"`
	Created customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	Updated customtypes.RFC3339TimeStringValue `tfsdk:"updated"`
}

func (model *VPCSubnetFilterModel) parseVPCSubnets(
	ctx context.Context,
	subnets []linodego.VPCSubnet,
) diag.Diagnostics {
	parseSubnet := func(subnet linodego.VPCSubnet) (VPCSubnetModel, diag.Diagnostics) {
		var s VPCSubnetModel
		s.ID = types.Int64Value(int64(subnet.ID))
		s.Label = types.StringValue(subnet.Label)
		s.IPv4 = types.StringValue(subnet.IPv4)

		linodes, diags := types.ListValueFrom(ctx, types.Int64Type, subnet.Linodes)
		if diags.HasError() {
			return s, diags
		}
		s.Linodes = linodes

		s.Created = customtypes.RFC3339TimeStringValue{
			StringValue: helper.NullableTimeToFramework(subnet.Created),
		}
		s.Updated = customtypes.RFC3339TimeStringValue{
			StringValue: helper.NullableTimeToFramework(subnet.Updated),
		}

		return s, nil
	}

	result := make([]VPCSubnetModel, len(subnets))

	for i, s := range subnets {
		subnet, diags := parseSubnet(s)
		if diags.HasError() {
			return diags
		}
		result[i] = subnet
	}

	model.VPCSubnets = result

	return nil
}
