package vpc

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

type VPCResourceModel struct {
	ID                   types.Int64                        `tfsdk:"id"`
	Label                types.String                       `tfsdk:"label"`
	Description          types.String                       `tfsdk:"description"`
	Region               types.String                       `tfsdk:"region"`
	Subnets              types.List                         `tfsdk:"subnets"`
	Created              customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	Updated              customtypes.RFC3339TimeStringValue `tfsdk:"updated"`
	SubnetsCreateOptions []VPCSubnetCreateOpts              `tfsdk:"subnets_create_options"`
}

type VPCDataSourceModel struct {
	ID          types.Int64                        `tfsdk:"id"`
	Label       types.String                       `tfsdk:"label"`
	Description types.String                       `tfsdk:"description"`
	Region      types.String                       `tfsdk:"region"`
	Subnets     types.List                         `tfsdk:"subnets"`
	Created     customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	Updated     customtypes.RFC3339TimeStringValue `tfsdk:"updated"`
}

type VPCSubnetCreateOpts struct {
	Label types.String `tfsdk:"label"`
	IPv4  types.String `tfsdk:"ipv4"`
}

func (d *VPCResourceModel) parseComputedAttributes(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(vpc.ID))
	d.Description = types.StringValue(vpc.Description)

	if vpc.Created != nil {
		d.Created = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringValue(vpc.Created.Format(time.RFC3339)),
		}
	} else {
		d.Created = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringNull(),
		}
	}

	if vpc.Updated != nil {
		d.Updated = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringValue(vpc.Updated.Format(time.RFC3339)),
		}
	} else {
		d.Created = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringNull(),
		}
	}

	subnetList, diags := parseSubnets(ctx, vpc.Subnets)
	if diags.HasError() {
		return diags
	}

	d.Subnets = *subnetList

	return nil
}

func (d *VPCResourceModel) parseVPC(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.Label = types.StringValue(vpc.Label)
	d.Region = types.StringValue(vpc.Region)

	subnetCreateOpts := make([]VPCSubnetCreateOpts, len(vpc.Subnets))

	for i, s := range vpc.Subnets {
		var createOpts VPCSubnetCreateOpts
		createOpts.IPv4 = types.StringValue(s.IPv4)
		createOpts.Label = types.StringValue(s.Label)
		subnetCreateOpts[i] = createOpts
	}

	d.SubnetsCreateOptions = subnetCreateOpts

	return d.parseComputedAttributes(ctx, vpc)
}

func (d *VPCDataSourceModel) parseVPC(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.Label = types.StringValue(vpc.Label)
	d.Description = types.StringValue(vpc.Description)
	d.Region = types.StringValue(vpc.Region)
	d.ID = types.Int64Value(int64(vpc.ID))

	if vpc.Created != nil {
		d.Created = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringValue(vpc.Created.Format(time.RFC3339)),
		}
	}

	if vpc.Updated != nil {
		d.Updated = customtypes.RFC3339TimeStringValue{
			StringValue: types.StringValue(vpc.Updated.Format(time.RFC3339)),
		}
	}

	subnetList, diags := parseSubnets(ctx, vpc.Subnets)
	if diags.HasError() {
		return diags
	}

	d.Subnets = *subnetList
	return nil
}

func parseSubnets(
	ctx context.Context,
	vpcSubnets []linodego.VPCSubnet,
) (*basetypes.ListValue, diag.Diagnostics) {
	subnets := make([]attr.Value, len(vpcSubnets))

	for i, subnet := range vpcSubnets {
		s := make(map[string]attr.Value)

		s["id"] = types.Int64Value(int64(subnet.ID))
		s["label"] = types.StringValue(subnet.Label)
		s["ipv4"] = types.StringValue(subnet.IPv4)

		if subnet.Created != nil {
			s["created"] = customtypes.RFC3339TimeStringValue{
				StringValue: types.StringValue(subnet.Created.Format(time.RFC3339)),
			}
		}

		if subnet.Updated != nil {
			s["updated"] = customtypes.RFC3339TimeStringValue{
				StringValue: types.StringValue(subnet.Updated.Format(time.RFC3339)),
			}
		}

		linodes, diags := types.ListValueFrom(ctx, types.Int64Type, subnet.Linodes)
		if diags.HasError() {
			return nil, diags
		}
		s["linodes"] = linodes

		obj, diags := types.ObjectValue(subnetObjectType.AttrTypes, s)
		if diags.HasError() {
			return nil, diags
		}

		subnets[i] = obj
	}

	subnetList, diags := basetypes.NewListValue(
		subnetObjectType,
		subnets,
	)
	if diags.HasError() {
		return nil, diags
	}

	return &subnetList, nil
}
