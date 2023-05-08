package helper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// ExpandFrameworkSet expands a framework types.Set into a primitive Go slice
func ExpandFrameworkSet[T any](ctx context.Context, set types.Set) ([]T, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var setValue basetypes.SetValue

	diagnostics.Append(set.ElementsAs(ctx, &setValue, false)...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	var result []T

	setValue.ElementsAs(ctx, &result, false)

	return result, diagnostics
}
