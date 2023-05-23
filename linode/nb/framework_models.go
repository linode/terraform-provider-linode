package nb

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

func (data *DataSourceModel) parseNodeBalancer(
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

	tags, diag := types.SetValueFrom(ctx, types.StringType, nodebalancer.Tags)
	if diag.HasError() {
		return diag
	}
	data.Tags = tags

	transfer, diag := parseTransfer(ctx, nodebalancer.Transfer)
	if diag.HasError() {
		return diag
	}
	data.Transfer = *transfer

	return nil
}

func parseTransfer(
	ctx context.Context,
	transfer linodego.NodeBalancerTransfer,
) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["in"] = toFloat64ValueWithDefault(transfer.In)
	result["out"] = toFloat64ValueWithDefault(transfer.Out)
	result["total"] = toFloat64ValueWithDefault(transfer.Total)

	transferObj, diag := types.ObjectValue(transferObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	resultList, diag := basetypes.NewListValue(
		transferObjectType,
		[]attr.Value{transferObj},
	)
	if diag.HasError() {
		return nil, diag
	}
	return &resultList, nil
}

// returns a Float64 with default value 0 if nil or a known value.
func toFloat64ValueWithDefault(value *float64) basetypes.Float64Value {
	if value != nil {
		return types.Float64PointerValue(value)
	} else {
		return types.Float64Value(0)
	}
}
