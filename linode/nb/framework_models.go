package nb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func (data *DataSourceModel) ParseNodeBalancer(
	ctx context.Context,
	nodebalancer *linodego.NodeBalancer,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(nodebalancer.ID))
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

	transferObj, diags := types.ObjectValue(transferObjectType.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}

	resultList, diags := basetypes.NewListValue(
		transferObjectType,
		[]attr.Value{transferObj},
	)
	if diags.HasError() {
		return nil, diags
	}
	return &resultList, nil
}
