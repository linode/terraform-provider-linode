package helper

import (
	"fmt"
	"math"
	"reflect"
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

func StringAliasSliceToStringSlice[T any](obj []T) ([]string, error) {
	var result []string

	var typeOfT T
	if reflect.TypeOf(typeOfT).Kind() != reflect.String {
		return nil, fmt.Errorf("type T is not an alias of string")
	}

	for _, v := range obj {
		strValue := reflect.ValueOf(v).String()
		result = append(result, strValue)
	}

	return result, nil
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
