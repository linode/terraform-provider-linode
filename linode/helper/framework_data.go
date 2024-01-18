package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func KeepOrUpdateStringPointer(original types.String, updated *string, preserveKnown bool) types.String {
	return KeepOrUpdateValue(original, types.StringPointerValue(updated), preserveKnown)
}

func KeepOrUpdateInt64Pointer(original types.Int64, updated *int64, preserveKnown bool) types.Int64 {
	return KeepOrUpdateValue(original, types.Int64PointerValue(updated), preserveKnown)
}

func KeepOrUpdateIntPointer(original types.Int64, updated *int, preserveKnown bool) types.Int64 {
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
