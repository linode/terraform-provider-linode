package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

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
func GetBaseSafeFwValueConverter[T any, U attr.Value](c SafeFwValueConverter[T, U]) SafeFwValueConverter[T, attr.Value] {
	return func(v T) attr.Value {
		return c(v)
	}
}
