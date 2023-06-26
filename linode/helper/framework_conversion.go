package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetValueIfNotNull - assign StringNull() safely without throwing error.
// e.g. new value: .rev_note: was null, but now cty.StringVal("")
func GetValueIfNotNull(val string) basetypes.StringValue {
	res := types.StringValue(val)

	if res == types.StringValue("") {
		res = types.StringNull()
	}

	return res
}

// GetStringPtrWithDefault returns a types.StringValue if the given pointer is
// not null, else it returns the provided default value.
func GetStringPtrWithDefault(val *string, def string) types.String {
	if val != nil {
		return types.StringValue(*val)
	}

	return types.StringValue(def)
}

// StringSliceToFramework converts the given string slice
// into a framework-compatible slice of types.String.
func StringSliceToFramework(val []string) []types.String {
	if val == nil {
		return nil
	}

	result := make([]types.String, len(val))

	for i, v := range val {
		result[i] = types.StringValue(v)
	}

	return result
}

// FrameworkSliceToString converts the given types.String slice
// into a primitive string slice.
func FrameworkSliceToString(val []types.String) []string {
	if val == nil {
		return nil
	}

	result := make([]string, len(val))

	for i, v := range val {
		result[i] = v.ValueString()
	}

	return result
}
