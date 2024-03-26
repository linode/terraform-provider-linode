package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// StringSliceToFrameworkValueSlice converts the given string slice
// into a framework-compatible slice of attr.Value.
func StringSliceToFrameworkValueSlice(val []string) []attr.Value {
	return GenericSliceToFramework(val, func(v string) attr.Value {
		return types.StringValue(v)
	})
}

type (
	// A function that converts something into a framework-compatible attr.Value
	// with a diag.Diagnostics for recording any error occurred.
	FwValueConverter[T any, U attr.Value] func(T) (U, diag.Diagnostics)

	// A function that converts something into a framework-compatible attr.Value
	SafeFwValueConverter[T any, U attr.Value] func(T) U
)

// Returns a converter that echos the framework value without any conversion
func FwValueEchoConverter() FwValueConverter[attr.Value, attr.Value] {
	return func(v attr.Value) (attr.Value, diag.Diagnostics) {
		return v, nil
	}
}

// Returns a converter that echos the framework value without any
// conversion without diags.Diagnostics support
func SafeFwValueEchoConverter() SafeFwValueConverter[attr.Value, attr.Value] {
	return func(v attr.Value) attr.Value { return v }
}

// Returns a converter that cast the result of the passed in converter
// to be base framework value (aka attr.Value)
func GetBaseFwValueConverter[T any, U attr.Value](c FwValueConverter[T, U]) FwValueConverter[T, attr.Value] {
	return func(v T) (attr.Value, diag.Diagnostics) {
		return c(v)
	}
}

// Returns a converter that cast the result of the passed in converter to
// be base framework value (aka attr.Value) without diags.Diagnostics support
func GetBaseSafeFwValueConverter[T any, U attr.Value](c FwValueConverter[T, U]) FwValueConverter[T, attr.Value] {
	return func(v T) (attr.Value, diag.Diagnostics) {
		return c(v)
	}
}

// GenericSliceToFramework converts the given generic slice
// into a framework-compatible slice of attr.Value.
func GenericSliceToFramework[T any, U attr.Value](
	val []T, converter SafeFwValueConverter[T, U],
) []U {
	if val == nil {
		return nil
	}

	result := make([]U, len(val))

	for i, v := range val {
		resultVal := converter(v)
		result[i] = resultVal
	}

	return result
}

// GenericSliceToFrameworkWithDiags converts the given generic slice
// into a framework-compatible slice of attr.Value with a converter and
// diag.Diagnostics for recording any error occurred.
func GenericSliceToFrameworkWithDiags[T any, U attr.Value](
	val []T, converter FwValueConverter[T, U], diags *diag.Diagnostics,
) []U {
	if val == nil {
		return nil
	}

	result := make([]U, len(val))

	for i, v := range val {
		resultVal, d := converter(v)
		if d.HasError() {
			diags.Append(d...)
			return nil
		}
		result[i] = resultVal
	}

	return result
}

// GenericSliceToList converts the given generic slice
// into a framework-compatible value of types.List with a FwValueConverter.
func GenericSliceToList[T any, V attr.Value](
	val []T, elementType attr.Type, converter FwValueConverter[T, V], diags *diag.Diagnostics,
) types.List {
	if val == nil {
		return types.ListNull(elementType)
	}

	result := GenericSliceToFrameworkWithDiags(val, GetBaseFwValueConverter(converter), diags)
	if diags.HasError() {
		return types.ListNull(elementType)
	}

	listResult, newDiags := types.ListValue(elementType, result)
	diags.Append(newDiags...)

	return listResult
}

// FrameworkSliceToString converts the given Framework slice
// into a slice of strings.
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

// IntSliceToFramework converts the given int slice
// into a framework-compatible slice of types.String.
func IntSliceToFramework(val []int) []types.Int64 {
	return GenericSliceToFramework(val, func(v int) types.Int64 {
		return types.Int64Value(int64(v))
	})
}
