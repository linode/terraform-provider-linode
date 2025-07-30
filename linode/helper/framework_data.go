package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// KeepOrUpdateValue is a generic function to keep the original value if it is known when preserveKnown is true,
// or update it otherwise
func KeepOrUpdateValue[T attr.Value](original T, updated T, preserveKnown bool) T {
	if preserveKnown && !original.IsUnknown() {
		return original
	}
	return updated
}

func KeepOrUpdateNestedObject[T any](
	ctx context.Context,
	original types.Object,
	preserveKnown bool,
	diags *diag.Diagnostics,
	flatten func(*T, bool, *diag.Diagnostics),
) *types.Object {
	return KeepOrUpdateNestedObjectWithTypes(
		ctx, original, original.AttributeTypes(ctx), preserveKnown, diags, flatten,
	)
}

// This function is necessary when explicit attributes are needed for flatten the `original`
// nested object.
//
// In some cases `original` won't contain the type of its attributes. For example, a
// double nested object (nested object in another nested object) in a model; when the
// parent nested object is null or unknown, `object.As` won't put the attributes into
// the child nested object. Passing explicit attributeTypes will then be necessary.
//
// Checkout the corresponding unit tests for more details.
func KeepOrUpdateNestedObjectWithTypes[T any](
	ctx context.Context,
	original types.Object,
	attributeTypes map[string]attr.Type,
	preserveKnown bool,
	diags *diag.Diagnostics,
	flatten func(*T, bool, *diag.Diagnostics),
) *types.Object {
	if preserveKnown && original.IsNull() {
		return &original
	}

	var attrModel T

	if !original.IsUnknown() && !original.IsNull() {
		diags.Append(
			original.As(ctx, &attrModel, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})...,
		)
		if diags.HasError() {
			return nil
		}
	}

	preserveKnown = preserveKnown && !original.IsUnknown()

	flatten(&attrModel, preserveKnown, diags)

	updated, newDiags := types.ObjectValueFrom(
		ctx,
		attributeTypes,
		attrModel,
	)
	diags.Append(newDiags...)
	if diags.HasError() {
		return nil
	}

	return &updated
}
