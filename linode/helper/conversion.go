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

func SafeInt64ToInt(number int64, diags *diag.Diagnostics) int {
	if number > math.MaxInt32 {
		diags.AddError(
			"Failed int64 to int conversion",
			"Integer %v is larger than the upper bound of int32",
		)
	} else if number < math.MinInt32 {
		diags.AddError(
			"Failed int64 to int conversion",
			"Integer %v is smaller than the lower bound of int32",
		)
	}
	return int(number)
}
