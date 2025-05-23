package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func KeepOrUpdateString(original types.String, updated string, preserveKnown bool) types.String {
	return KeepOrUpdateValue(original, types.StringValue(updated), preserveKnown)
}

func KeepOrUpdateInt64(original types.Int64, updated int64, preserveKnown bool) types.Int64 {
	return KeepOrUpdateValue(original, types.Int64Value(updated), preserveKnown)
}

func KeepOrUpdateBool(original types.Bool, updated bool, preserveKnown bool) types.Bool {
	return KeepOrUpdateValue(original, types.BoolValue(updated), preserveKnown)
}

func KeepOrUpdateStringSet(
	original types.Set, updated []string, preserveKnown bool, diags *diag.Diagnostics,
) types.Set {
	return KeepOrUpdateSet(
		types.StringType, original, StringSliceToFrameworkValueSlice(updated), preserveKnown, diags,
	)
}

func KeepOrUpdateIntSet(
	original types.Set, updated []int, preserveKnown bool, diags *diag.Diagnostics,
) types.Set {
	return KeepOrUpdateSet(
		types.Int64Type, original, IntSliceToFrameworkValueSlice(updated), preserveKnown, diags,
	)
}

func KeepOrUpdateStringMap(
	ctx context.Context,
	original types.Map,
	updated map[string]string,
	preserveKnown bool,
	diags *diag.Diagnostics,
) types.Map {
	mapValue, newDiags := types.MapValueFrom(ctx, types.StringType, updated)
	diags.Append(newDiags...)

	if diags.HasError() {
		return mapValue
	}

	return KeepOrUpdateValue(original, mapValue, preserveKnown)
}

func KeepOrUpdateSet(
	elementType attr.Type, original types.Set, updated []attr.Value, preserveKnown bool, diags *diag.Diagnostics,
) types.Set {
	setValue, newDiags := types.SetValue(elementType, updated)
	diags.Append(newDiags...)

	if diags.HasError() {
		return setValue
	}

	return KeepOrUpdateValue(original, setValue, preserveKnown)
}

func KeepOrUpdateStringPointer(original types.String, updated *string, preserveKnown bool) types.String {
	return KeepOrUpdateValue(original, types.StringPointerValue(updated), preserveKnown)
}

func KeepOrUpdateInt64Pointer(original types.Int64, updated *int64, preserveKnown bool) types.Int64 {
	return KeepOrUpdateValue(original, types.Int64PointerValue(updated), preserveKnown)
}

func KeepOrUpdateInt32Pointer(original types.Int32, updated *int32, preserveKnown bool) types.Int32 {
	return KeepOrUpdateValue(original, types.Int32PointerValue(updated), preserveKnown)
}

func KeepOrUpdateFloat64Pointer(original types.Float64, updated *float64, preserveKnown bool) types.Float64 {
	return KeepOrUpdateValue(original, types.Float64PointerValue(updated), preserveKnown)
}

func KeepOrUpdateIntPointer(original types.Int64, updated *int, preserveKnown bool) types.Int64 {
	// There is not a built in function in `types` library of the framework.
	// Manually handle it here
	if updated == nil {
		return KeepOrUpdateValue(original, types.Int64Null(), preserveKnown)
	}
	return KeepOrUpdateInt64(original, int64(*updated), preserveKnown)
}

func KeepOrUpdateBoolPointer(original types.Bool, updated *bool, preserveKnown bool) types.Bool {
	return KeepOrUpdateValue(original, types.BoolPointerValue(updated), preserveKnown)
}

func KeepOrUpdateValue[T attr.Value](original T, updated T, preserveKnown bool) T {
	if preserveKnown && !original.IsUnknown() {
		return original
	}
	return updated
}
