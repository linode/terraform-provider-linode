package nb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

// NodebalancerModel describes the Terraform resource data model to match the
// resource schema.
type NodebalancerModel struct {
	ID                 types.Int64                        `tfsdk:"id"`
	Label              types.String                       `tfsdk:"label"`
	Region             types.String                       `tfsdk:"region"`
	ClientConnThrottle types.Int64                        `tfsdk:"client_conn_throttle"`
	Hostname           types.String                       `tfsdk:"hostname"`
	Ipv4               types.String                       `tfsdk:"ipv4"`
	Ipv6               types.String                       `tfsdk:"ipv6"`
	Created            customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	Updated            customtypes.RFC3339TimeStringValue `tfsdk:"updated"`
	Transfer           types.List                         `tfsdk:"transfer"`
	Tags               types.Set                          `tfsdk:"tags"`
}

type nbModelV0 struct {
	ID                 types.Int64  `tfsdk:"id"`
	Label              types.String `tfsdk:"label"`
	Region             types.String `tfsdk:"region"`
	ClientConnThrottle types.Int64  `tfsdk:"client_conn_throttle"`
	Hostname           types.String `tfsdk:"hostname"`
	Ipv4               types.String `tfsdk:"ipv4"`
	Ipv6               types.String `tfsdk:"ipv6"`
	Created            types.String `tfsdk:"created"`
	Updated            types.String `tfsdk:"updated"`
	Tags               types.Set    `tfsdk:"tags"`
	Transfer           types.Map    `tfsdk:"transfer"`
}

func (data *NodebalancerModel) ParseNonComputedAttrs(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Label = types.StringPointerValue(nodebalancer.Label)

	tags, diags := types.SetValueFrom(ctx, types.StringType, helper.StringSliceToFramework(nodebalancer.Tags))
	if diags.HasError() {
		return diags
	}
	data.Tags = tags

	return nil
}

func (data *NodebalancerModel) ParseComputedAttrs(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(nodebalancer.ID))
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.Ipv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.Ipv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = customtypes.RFC3339TimeStringValue{
		StringValue: types.StringValue(nodebalancer.Created.Format(time.RFC3339)),
	}
	data.Updated = customtypes.RFC3339TimeStringValue{
		StringValue: types.StringValue(nodebalancer.Updated.Format(time.RFC3339)),
	}

	transfer, diags := parseTransfer(ctx, nodebalancer.Transfer)
	if diags.HasError() {
		return diags
	}

	data.Transfer = *transfer

	return nil
}

func parseTransfer(
	ctx context.Context,
	transfer linodego.NodeBalancerTransfer,
) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["in"] = helper.Float64PointerValueWithDefault(transfer.In)
	result["out"] = helper.Float64PointerValueWithDefault(transfer.Out)
	result["total"] = helper.Float64PointerValueWithDefault(transfer.Total)

	transferObj, diags := types.ObjectValue(TransferObjectType.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}

	resultList, diags := types.ListValueFrom(ctx, TransferObjectType, []attr.Value{transferObj})

	if diags.HasError() {
		return nil, diags
	}

	return &resultList, nil
}

func UpgradeResourceStateValue(val string) (basetypes.Float64Value, diag.Diagnostic) {
	if val == "" {
		return types.Float64Value(0), nil
	}
	result, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return types.Float64Null(), diag.NewErrorDiagnostic(
			"Failed to upgrade state.",
			err.Error(),
		)
	}
	return types.Float64Value(result), nil
}
