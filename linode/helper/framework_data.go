package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func KeepOrUpdateString(original types.String, updated string, preserveKnown bool) types.String {
	return KeepOrUpdateStringValue(original, types.StringValue(updated), preserveKnown)
}

func KeepOrUpdateInt64(original types.Int64, updated int64, preserveKnown bool) types.Int64 {
	return KeepOrUpdateInt64Value(original, types.Int64Value(updated), preserveKnown)
}

func KeepOrUpdateBool(original types.Bool, updated bool, preserveKnown bool) types.Bool {
	return KeepOrUpdateBoolValue(original, types.BoolValue(updated), preserveKnown)
}

func KeepOrUpdateStringValue(original types.String, updated types.String, preserveKnown bool) types.String {
	return KeepOrUpdateValue(original, updated, preserveKnown).(types.String)
}

func KeepOrUpdateInt64Value(original types.Int64, updated types.Int64, preserveKnown bool) types.Int64 {
	return KeepOrUpdateValue(original, updated, preserveKnown).(types.Int64)
}

func KeepOrUpdateBoolValue(original types.Bool, updated types.Bool, preserveKnown bool) types.Bool {
	return KeepOrUpdateValue(original, updated, preserveKnown).(types.Bool)
}

func KeepOrUpdateListValue(
	original types.List,
	updated types.List,
	elementType attr.Type,
	preserveKnown bool,
) types.List {
	return KeepOrUpdateValue(original, updated, preserveKnown).(types.List)
}

func KeepOrUpdateKeepSetValue(
	original types.Set,
	updated types.Set,
	elementType attr.Type,
	preserveKnown bool,
) types.Set {
	return KeepOrUpdateValue(original, updated, preserveKnown).(types.Set)
}

func KeepOrUpdateValue(original attr.Value, updated attr.Value, preserveKnown bool) attr.Value {
	if preserveKnown && !original.IsUnknown() {
		return original
	}
	return updated
}
