package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// returns a Float64Value with default value 0 if nil or a known value.
func Float64PointerValueWithDefault(value *float64) types.Float64 {
	if value != nil {
		return types.Float64PointerValue(value)
	} else {
		return types.Float64Value(0)
	}
}

// returns an Int64Value with default value 0 if nil or a known value.
func IntPointerValueWithDefault(value *int) types.Int64 {
	if value != nil {
		return types.Int64Value(int64(*value))
	} else {
		return types.Int64Value(0)
	}
}

// returns an StringValue with default value "" if nil or a known value.
func StringPointerValueWithDefault(value *string) types.String {
	if value != nil {
		return types.StringValue(*value)
	} else {
		return types.StringValue("")
	}
}

func ValueStringPointerWithNilForUnknownAndEmptyString(stringValue types.String) *string {
	if stringValue.ValueString() == "" {
		return nil
	} else {
		return stringValue.ValueStringPointer()
	}
}
