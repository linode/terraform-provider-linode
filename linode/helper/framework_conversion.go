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

// FrameworkToStringSlice converts given []types.String to []string
func FrameworkToStringSlice(vals []types.String) []string {
	if vals == nil {
		return nil
	}

	result := make([]string, len(vals))

	for i, v := range vals {
		result[i] = v.ValueString()
	}

	return result
}
