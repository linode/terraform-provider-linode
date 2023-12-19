package vpcsubnets

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnet"
)

// VPCSubnetFilterModel describes the Terraform resource data model to match the
// resource schema.
type VPCSubnetFilterModel struct {
	ID         types.String                     `tfsdk:"id"`
	VPCId      types.Int64                      `tfsdk:"vpc_id"`
	Filters    frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCSubnets []VPCSubnetModel                 `tfsdk:"vpc_subnets"`
}

type VPCSubnetModel struct {
	ID      types.Int64       `tfsdk:"id"`
	Label   types.String      `tfsdk:"label"`
	IPv4    types.String      `tfsdk:"ipv4"`
	Linodes types.List        `tfsdk:"linodes"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	Updated timetypes.RFC3339 `tfsdk:"updated"`
}

func (model *VPCSubnetFilterModel) FlattenSubnets(
	ctx context.Context,
	subnets []linodego.VPCSubnet,
	preserveKnown bool,
) diag.Diagnostics {
	parseSubnet := func(subnet linodego.VPCSubnet) (VPCSubnetModel, diag.Diagnostics) {
		var s VPCSubnetModel
		s.ID = helper.KeepOrUpdateInt64(s.ID, int64(subnet.ID), preserveKnown)
		s.Label = helper.KeepOrUpdateString(s.Label, subnet.Label, preserveKnown)
		s.IPv4 = helper.KeepOrUpdateString(s.IPv4, subnet.IPv4, preserveKnown)

		linodes := make([]types.Object, len(subnet.Linodes))

		for i, inst := range subnet.Linodes {
			linodeObj, d := vpcsubnet.FlattenSubnetLinode(ctx, inst)
			if d.HasError() {
				return s, d
			}

			linodes[i] = *linodeObj
		}

		linodesList, d := types.ListValueFrom(ctx, vpcsubnet.LinodeObjectType, linodes)
		if d.HasError() {
			return s, d
		}
		s.Linodes = helper.KeepOrUpdateValue(s.Linodes, linodesList, preserveKnown)

		s.Created = helper.KeepOrUpdateValue(
			s.Created,
			timetypes.NewRFC3339TimePointerValue(subnet.Created),
			preserveKnown,
		)
		s.Updated = helper.KeepOrUpdateValue(
			s.Updated,
			timetypes.NewRFC3339TimePointerValue(subnet.Updated),
			preserveKnown,
		)

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
