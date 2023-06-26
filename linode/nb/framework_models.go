package nb

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type TransferModelEntry struct {
	In    types.Float64 `tfsdk:"in"`
	Out   types.Float64 `tfsdk:"out"`
	Total types.Float64 `tfsdk:"total"`
}

// NodebalancerModel describes the Terraform resource data model to match the
// resource schema.
type NodebalancerModel struct {
	ID                 types.Int64          `tfsdk:"id"`
	Label              types.String         `tfsdk:"label"`
	Region             types.String         `tfsdk:"region"`
	ClientConnThrottle types.Int64          `tfsdk:"client_conn_throttle"`
	Hostname           types.String         `tfsdk:"hostname"`
	Ipv4               types.String         `tfsdk:"ipv4"`
	Ipv6               types.String         `tfsdk:"ipv6"`
	Created            types.String         `tfsdk:"created"`
	Updated            types.String         `tfsdk:"updated"`
	Transfer           []TransferModelEntry `tfsdk:"transfer"`
	Tags               types.Set            `tfsdk:"tags"`
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

func (data *NodebalancerModel) parseNodeBalancer(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.Label = types.StringPointerValue(nodebalancer.Label)
	data.Region = types.StringValue(nodebalancer.Region)
	data.ClientConnThrottle = types.Int64Value(int64(nodebalancer.ClientConnThrottle))
	data.Hostname = types.StringPointerValue(nodebalancer.Hostname)
	data.Ipv4 = types.StringPointerValue(nodebalancer.IPv4)
	data.Ipv6 = types.StringPointerValue(nodebalancer.IPv6)
	data.Created = types.StringValue(nodebalancer.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(nodebalancer.Updated.Format(time.RFC3339))

	tags, diags := types.SetValueFrom(ctx, types.StringType, nodebalancer.Tags)
	if diags.HasError() {
		return diags
	}
	data.Tags = tags
	data.Transfer = parseTransfer(nodebalancer.Transfer)

	return nil
}

func parseTransfer(
	transfer linodego.NodeBalancerTransfer,
) []TransferModelEntry {
	var entry TransferModelEntry

	entry.In = helper.Float64PointerValueWithDefault(transfer.In)
	entry.Out = helper.Float64PointerValueWithDefault(transfer.Out)
	entry.Total = helper.Float64PointerValueWithDefault(transfer.Total)

	return []TransferModelEntry{entry}
}

func stringToFloat64Value(val string) (basetypes.Float64Value, diag.Diagnostic) {
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
