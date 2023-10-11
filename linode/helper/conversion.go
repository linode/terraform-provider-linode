package helper

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TypedSliceToAny[T any](obj []T) []any {
	result := make([]any, len(obj))

	for i, v := range obj {
		result[i] = v
	}

	return result
}

func AnySliceToTyped[T any](obj []any) []T {
	result := make([]T, len(obj))

	for i, v := range obj {
		result[i] = v.(T)
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

func SafeInt64ToInt(number int64) (int, error) {
	if number > math.MaxInt || number < math.MinInt {
		return 0, fmt.Errorf("int64 value %v is out of range for int", number)
	}
	return int(number), nil
}

func SafeFloat64ToInt(number float64) (int, error) {
	if number > float64(math.MaxInt) || number < float64(math.MinInt) {
		return 0, fmt.Errorf("float64 value %v is out of range for int64", number)
	}
	return int(number), nil
}
