package helper

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

func TypedSliceToAny[T any](obj []T) []any {
	return MapSlice(obj, func(v T) any {
		return v
	})
}

func AnySliceToTyped[T any](obj []any) []T {
	return MapSlice(obj, func(v any) T {
		return v.(T)
	})
}

func StringTypedMapToAny[T any](m map[string]T) map[string]any {
	return MapMap(m, func(k string, v T) (string, any) {
		return k, v
	})
}

func StringAnyMapToTyped[T any](m map[string]any) map[string]T {
	return MapMap(m, func(k string, v any) (string, T) {
		return k, v.(T)
	})
}

func StringAliasSliceToStringSlice[T ~string](obj []T) []string {
	var result []string

	for _, v := range obj {
		strValue := reflect.ValueOf(v).String()
		result = append(result, strValue)
	}

	return result
}

func StringToInt64(s string, diags *diag.Diagnostics) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Invalid number string: %v", s),
			err.Error(),
		)
	}
	return num
}

func StringToInt(s string, diags *diag.Diagnostics) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Invalid number string: %v", s),
			err.Error(),
		)
	}
	return num
}

func FrameworkSafeInt64ToInt(number int64, diags *diag.Diagnostics) int {
	result, err := SafeInt64ToInt(number)
	if err != nil {
		diags.AddError(
			"Failed int64 to int conversion",
			err.Error(),
		)
	}
	return result
}

func FrameworkSafeInt64PointerToIntPointer(number *int64, diags *diag.Diagnostics) *int {
	if number == nil {
		return nil
	}

	result, err := SafeInt64ToInt(*number)
	if err != nil {
		diags.AddError(
			"Failed int64 pointer to int pointer conversion",
			err.Error(),
		)
	}

	return linodego.Pointer(result)
}

func FrameworkSafeInt64ValueToIntDoublePointerWithUnknownToNil(v types.Int64, diags *diag.Diagnostics) **int {
	if v.IsUnknown() {
		return linodego.DoublePointerNull[int]()
	}

	return linodego.Pointer(FrameworkSafeInt64PointerToIntPointer(v.ValueInt64Pointer(), diags))
}

func FrameworkSafeInt64ValueToIntPointerWithUnknownToNil(v types.Int64, diags *diag.Diagnostics) *int {
	if v.IsUnknown() {
		return nil
	}

	return linodego.Pointer(FrameworkSafeInt64ToInt(v.ValueInt64(), diags))
}

func ValueBoolPointerWithUnknownToNil(v types.Bool) *bool {
	if v.IsUnknown() {
		return nil
	}

	return v.ValueBoolPointer()
}

func ValueStringPointerWithUnknownToNil(v types.String) *string {
	if v.IsUnknown() {
		return nil
	}

	return v.ValueStringPointer()
}

func FrameworkSafeFloat64ToInt(number float64, diags *diag.Diagnostics) int {
	result, err := SafeFloat64ToInt(number)
	if err != nil {
		diags.AddError(
			"Failed float64 to int conversion",
			err.Error(),
		)
	}
	return result
}

func IntPtrToInt64Ptr(ptr *int) *int64 {
	if ptr == nil {
		return nil
	}
	val := int64(*ptr)
	return &val
}

func SafeInt64ToInt(number int64) (int, error) {
	if number > math.MaxInt || number < math.MinInt {
		return 0, fmt.Errorf("int64 value %v is out of range for int", number)
	}
	return int(number), nil
}

func SafeIntToInt32(number int) (int32, error) {
	if number > math.MaxInt32 || number < math.MinInt32 {
		return 0, fmt.Errorf("int value %v is out of range for int32", number)
	}
	return int32(number), nil
}

func SafeFloat64ToInt(number float64) (int, error) {
	if number > float64(math.MaxInt) || number < float64(math.MinInt) {
		return 0, fmt.Errorf("float64 value %v is out of range for int64", number)
	}
	return int(number), nil
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func FrameworkSafeStringToInt(val string, d *diag.Diagnostics) int {
	result, err := strconv.Atoi(val)
	if err != nil {
		d.Append(diag.NewErrorDiagnostic(
			"Failed to convert string to int",
			err.Error(),
		))
		return 0
	}

	return result
}
