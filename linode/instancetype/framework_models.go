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
	ID           types.String `tfsdk:"id"`
	Label        types.String `tfsdk:"label"`
	Disk         types.Int64  `tfsdk:"disk"`
	Class        types.String `tfsdk:"class"`
	Price        types.List   `tfsdk:"price"`
	Addons       types.List   `tfsdk:"addons"`
	RegionPrices types.List   `tfsdk:"region_prices"`
	NetworkOut   types.Int64  `tfsdk:"network_out"`
	Memory       types.Int64  `tfsdk:"memory"`
	Transfer     types.Int64  `tfsdk:"transfer"`
	VCPUs        types.Int64  `tfsdk:"vcpus"`
}

func (data *DataSourceModel) ParseLinodeType(
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

	regionPrices, d := FlattenRegionPrices(linodeType.RegionPrices)
	if d.HasError() {
		return d
	}
	data.RegionPrices = *regionPrices

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

	resultList, diag := types.ListValue(
		priceObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func FlattenRegionPrices(prices []linodego.LinodeRegionPrice) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make([]attr.Value, len(prices))

	for i, price := range prices {
		obj, d := types.ObjectValue(priceObjectType.AttrTypes, map[string]attr.Value{
			"id":      types.StringValue(price.ID),
			"hourly":  types.Float64Value(float64(price.Hourly)),
			"monthly": types.Float64Value(float64(price.Monthly)),
		})
		if d.HasError() {
			return nil, d
		}

		result[i] = obj
	}

	priceList, d := basetypes.NewListValue(
		priceObjectType,
		result,
	)
	return &priceList, d
}
