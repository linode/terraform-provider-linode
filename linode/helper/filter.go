package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"encoding/json"
)

// FilterTypeFunc is a function that takes in a filter name and value,
// and returns the value converted to the correct filter type.
type FilterTypeFunc func(filterName string, value string) (interface{}, error)

// FilterSchema should be referenced in a schema configuration in order to
// enable filter functionality
func FilterSchema(validFilters []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Description:  "The name of the attribute to filter on.",
					ValidateFunc: validation.StringInSlice(validFilters, false),
					Required:     true,
				},
				"values": {
					Type:        schema.TypeList,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "The value(s) to be used in the filter.",
					Required:    true,
				},
			},
		},
	}
}

// ConstructFilterString constructs a Linode filter JSON string from each filter element in the schema
func ConstructFilterString(d *schema.ResourceData, typeFunc FilterTypeFunc) (string, error) {
	filters := d.Get("filter").([]interface{})
	resultMap := make(map[string]interface{})

	if len(filters) < 1 {
		return "{}", nil
	}

	var rootFilter []interface{}

	for _, filter := range filters {
		filter := filter.(map[string]interface{})

		name := filter["name"].(string)
		values := filter["values"].([]interface{})

		subFilter := make([]interface{}, len(values))

		for i, value := range values {
			value, err := typeFunc(name, value.(string))
			if err != nil {
				return "", err
			}

			valueFilter := make(map[string]interface{})
			valueFilter[name] = value

			subFilter[i] = valueFilter
		}

		rootFilter = append(rootFilter, map[string]interface{}{
			"+or": subFilter,
		})
	}

	resultMap["+and"] = rootFilter

	result, err := json.Marshal(resultMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
