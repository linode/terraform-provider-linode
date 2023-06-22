package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
// into a framework-compatible slice of attr.Value.
// works for List and Set attributes
func StringSliceToFramework(val []string) []attr.Value {
	if val == nil {
		return nil
	}

	result := make([]attr.Value, len(val))

	for i, v := range val {
		result[i] = types.StringValue(v)
	}

	return result
}

func FrameworkSetToStringSlice(
	ctx context.Context,
	vals basetypes.SetValue,
) []string {
	if vals.IsNull() {
		return nil
	}

	result := make([]string, len(vals.Elements()))
	diag := vals.ElementsAs(ctx, result, false)
	if diag.HasError() {
		return nil
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
