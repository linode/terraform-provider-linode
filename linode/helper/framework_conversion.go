package helper

import (
	"context"

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
func StringSliceToFramework(vals []string) []types.String {
	if vals == nil {
		return nil
	}

	result := make([]types.String, len(vals))

	for i, v := range vals {
		result[i] = types.StringValue(v)
	}

	return result
}

// FrameworkSetToStringSlice converts SetValue to []string
func FrameworkSetToStringSlice(ctx context.Context, vals basetypes.SetValue) []string {
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

// StringSliceToFrameworkSet converts []string to basetypes.SetValue
func StringSliceToFrameworkSet(vals []string) basetypes.SetValue {
	result, _ := basetypes.NewSetValueFrom(context.Background(), types.StringType, vals)
	return result
}
