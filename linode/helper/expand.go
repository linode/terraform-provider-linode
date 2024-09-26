package helper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandStringList(list []interface{}) []string {
	slice := make([]string, 0, len(list))
	for _, s := range list {
		if val, ok := s.(string); ok && val != "" {
			slice = append(slice, val)
		}
	}
	return slice
}

func ExpandStringSet(set *schema.Set) []string {
	return ExpandStringList(set.List())
}

func ExpandObjectList(list []any) []map[string]any {
	slice := make([]map[string]any, 0, len(list))
	for _, s := range list {
		if val, ok := s.(map[string]any); ok {
			slice = append(slice, val)
		}
	}
	return slice
}

func ExpandObjectSet(set *schema.Set) []map[string]any {
	return ExpandObjectList(set.List())
}

func ExpandIntList(list []any) []int {
	slice := make([]int, 0, len(list))
	for _, n := range list {
		if val, ok := n.(int); ok {
			slice = append(slice, val)
		}
	}
	return slice
}

func ExpandIntSet(set *schema.Set) []int {
	return ExpandIntList(set.List())
}

func ExpandFwInt64Set(set types.Set, diags *diag.Diagnostics) (result []int) {
	elements := set.Elements()
	result = make([]int, len(elements))

	for i, v := range elements {
		num, ok := v.(types.Int64)
		if !ok {
			diags.AddError(
				"Value Conversion Failed",
				fmt.Sprintf("Failed to convert %v to int", v),
			)
			return
		}
		result[i] = FrameworkSafeInt64ToInt(num.ValueInt64(), diags)
		if diags.HasError() {
			return
		}
	}
	return
}
