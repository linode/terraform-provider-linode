package helper

import (
	"fmt"
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

func StringToInt64(s string, diags diag.Diagnostics) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Invalid number string: %v", s),
			err.Error(),
		)
		return 0
	}

	return num
}
