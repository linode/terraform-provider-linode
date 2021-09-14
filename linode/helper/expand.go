package helper

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func ExpandIntList(list []interface{}) []int {
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
