package vpcsubnets

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcsubnet"
)

// Model describes the Terraform resource data model to match the
// resource schema.
type Model struct {
	ID         types.String                     `tfsdk:"id"`
	VPCId      types.Int64                      `tfsdk:"vpc_id"`
	Filters    frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCSubnets []ModelVPCSubnet                 `tfsdk:"vpc_subnets"`
}

// ModelVPCSubnet defines the structure of each entry returned by this data source.
// NOTE: This is redefined because of divergence between this model and the
// singular data source's model. We should investigate using composition in the future.
type ModelVPCSubnet struct {
	ID      types.Int64       `tfsdk:"id"`
	Label   types.String      `tfsdk:"label"`
	IPv4    types.String      `tfsdk:"ipv4"`
	IPv6    types.List        `tfsdk:"ipv6"`
	Linodes types.List        `tfsdk:"linodes"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	Updated timetypes.RFC3339 `tfsdk:"updated"`
}

func (model *Model) FlattenSubnets(
	ctx context.Context,
	subnets []linodego.VPCSubnet,
	preserveKnown bool,
) (diags diag.Diagnostics) {
	result := helper.MapSlice(
		subnets,
		func(subnet linodego.VPCSubnet) ModelVPCSubnet {
			var s ModelVPCSubnet
			s.ID = helper.KeepOrUpdateInt64(s.ID, int64(subnet.ID), preserveKnown)
			s.Label = helper.KeepOrUpdateString(s.Label, subnet.Label, preserveKnown)
			s.IPv4 = helper.KeepOrUpdateString(s.IPv4, subnet.IPv4, preserveKnown)

			ipv6AddressModels := helper.MapSlice(
				subnet.IPv6,
				func(subnet linodego.VPCIPv6Range) vpcsubnet.ResourceModelIPv6 {
					return vpcsubnet.ResourceModelIPv6{
						Range: customtypes.LinodeAutoAllocRangeValue{StringValue: types.StringValue(subnet.Range)},
					}
				},
			)

			ipv6AddressesList, d := types.ListValueFrom(ctx, vpcsubnet.ResourceModelIPv6ObjectType, ipv6AddressModels)
			diags.Append(d...)
			if diags.HasError() {
				return s
			}

			s.IPv6 = helper.KeepOrUpdateValue(
				s.IPv6,
				ipv6AddressesList,
				preserveKnown,
			)

			linodes := make([]types.Object, len(subnet.Linodes))

			for i, inst := range subnet.Linodes {
				linodeObj, d := vpcsubnet.FlattenSubnetLinode(ctx, inst)
				diags.Append(d...)
				if diags.HasError() {
					return s
				}

				linodes[i] = *linodeObj
			}

			linodesList, d := types.ListValueFrom(ctx, vpcsubnet.LinodeObjectType, linodes)
			diags.Append(d...)
			if diags.HasError() {
				return s
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

			return s
		},
	)
	if diags.HasError() {
		return
	}

	model.VPCSubnets = result

	return
}
