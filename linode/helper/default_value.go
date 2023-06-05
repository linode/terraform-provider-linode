package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// returns a Float64Value with default value 0 if nil or a known value.
func Float64PointerValueWithDefault(value *float64) basetypes.Float64Value {
	if value != nil {
		return types.Float64PointerValue(value)
	} else {
		return types.Float64Value(0)
	}
}

// returns an Int64Value with default value 0 if nil or a known value.
func IntPointerValueWithDefault(value *int) basetypes.Int64Value {
	if value != nil {
		return types.Int64Value(int64(*value))
	} else {
		return types.Int64Value(0)
	}
}
