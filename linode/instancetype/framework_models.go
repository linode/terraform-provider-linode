package instancetype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Label      types.String `tfsdk:"label"`
	Disk       types.Int64  `tfsdk:"disk"`
	Class      types.String `tfsdk:"class"`
	Price      types.List   `tfsdk:"price"`
	Addons     types.List   `tfsdk:"addons"`
	NetworkOut types.Int64  `tfsdk:"network_out"`
	Memory     types.Int64  `tfsdk:"memory"`
	Transfer   types.Int64  `tfsdk:"transfer"`
	VCPUs      types.Int64  `tfsdk:"vcpus"`
}

func (data *DataSourceModel) parseLinodeType(
	ctx context.Context, linodeType *linodego.LinodeType,
) diag.Diagnostics {
	data.ID = types.StringValue(linodeType.ID)
	data.Disk = types.Int64Value(int64(linodeType.Disk))
	data.Class = types.StringValue(string(linodeType.Class))

	price, diags := FlattenPrice(ctx, *linodeType.Price)
	if diags.HasError() {
		return diags
	}
	data.Price = *price

	data.Label = types.StringValue(linodeType.Label)

	addons, diags := FlattenAddons(ctx, *linodeType.Addons)
	if diags.HasError() {
		return diags
	}
	data.Addons = *addons

	data.NetworkOut = types.Int64Value(int64(linodeType.NetworkOut))
	data.Memory = types.Int64Value(int64(linodeType.Memory))
	data.Transfer = types.Int64Value(int64(linodeType.Transfer))
	data.VCPUs = types.Int64Value(int64(linodeType.VCPUs))

	return nil
}

func FlattenAddons(ctx context.Context, backup linodego.LinodeAddons) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	backups, diag := FlattenBackups(ctx, *backup.Backups)
	if diag.HasError() {
		return nil, diag
	}

	result["backups"] = backups

	obj, diag := types.ObjectValue(addonsObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		addonsObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func FlattenBackups(ctx context.Context, backup linodego.LinodeBackupsAddon) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	price, diag := FlattenPrice(ctx, *backup.Price)
	if diag.HasError() {
		return nil, diag
	}

	result["price"] = price

	obj, diag := types.ObjectValue(backupsObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		backupsObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func FlattenPrice(ctx context.Context, price linodego.LinodePrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["hourly"] = types.Float64Value(float64(price.Hourly))
	result["monthly"] = types.Float64Value(float64(price.Monthly))

	obj, diag := types.ObjectValue(priceObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		priceObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}
