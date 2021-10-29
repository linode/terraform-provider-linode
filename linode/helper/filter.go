package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"golang.org/x/crypto/sha3"
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
				"match_by": {
					Type:        schema.TypeString,
					Description: "The type of comparison to use for this filter.",
					Optional:    true,
					Default:     "exact",
					ValidateFunc: validation.StringInSlice([]string{"exact", "substring", "sub", "re", "regex"},
						false),
				},
			},
		},
	}
}

// OrderBySchema should be referenced in a schema configuration in order to
// enable filter ordering functionality
func OrderBySchema(validFilters []string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(validFilters, false),
		Description:  "The attribute to order the results by.",
	}
}

func OrderSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "asc",
		ValidateFunc: validation.StringInSlice([]string{"asc", "desc"}, false),
		Description:  "The order in which results should be returned.",
	}
}

// GetFilterID creates a unique ID specific to the current filter data source
func GetFilterID(d *schema.ResourceData) (string, error) {
	idMap := map[string]interface{}{
		"filter":   d.Get("filter"),
		"order":    d.Get("order"),
		"order_by": d.Get("order_by"),
	}

	result, err := json.Marshal(idMap)
	if err != nil {
		return "", err
	}

	hash := sha3.Sum512(result)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
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
		matchBy := filter["match_by"].(string)

		// Defer this logic to the client
		if matchBy != "exact" {
			continue
		}

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

	if len(rootFilter) < 1 {
		return "{}", nil
	}

	resultMap["+and"] = rootFilter

	if orderBy, ok := d.GetOk("order_by"); ok {
		resultMap["+order_by"] = orderBy
		resultMap["+order"] = d.Get("order")
	}

	result, err := json.Marshal(resultMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// FilterResults filters the given results on the client-side filters present in the resource
func FilterResults(d *schema.ResourceData, items []interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, item := range items {
		item := item.(map[string]interface{})

		match, err := itemMatchesFilter(d, item)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}

func itemMatchesFilter(d *schema.ResourceData, item map[string]interface{}) (bool, error) {
	filters := d.Get("filter").([]interface{})

	for _, filter := range filters {
		filter := filter.(map[string]interface{})

		name := filter["name"].(string)
		values := filter["values"].([]interface{})
		matchBy := filter["match_by"].(string)

		if matchBy == "exact" {
			continue
		}

		itemValue, ok := item[name]
		if !ok {
			return false, fmt.Errorf("\"%v\" is not a valid attribute", name)
		}

		valid, err := validateFilter(matchBy, name, ExpandStringList(values), itemValue)
		if err != nil {
			return false, err
		}

		if !valid {
			return false, nil
		}
	}

	return true, nil
}

func validateFilter(matchBy, name string, values []string, itemValue interface{}) (bool, error) {
	// Filter recursively on lists (tags, etc.)
	if items, ok := itemValue.([]string); ok {
		for _, item := range items {
			valid, err := validateFilter(matchBy, name, values, item)
			if err != nil {
				return false, err
			}

			if valid {
				return true, nil
			}
		}

		return false, nil
	}

	// Only string attributes should be considered
	itemValueStr, ok := itemValue.(string)
	if !ok {
		return false, fmt.Errorf("\"%s\" is not a string", name)
	}

	switch matchBy {
	case "substring", "sub":
		return validateFilterSubstring(values, itemValueStr)
	case "re", "regex":
		return validateFilterRegex(values, itemValueStr)
	}

	return true, nil
}

func validateFilterSubstring(values []string, result string) (bool, error) {
	for _, value := range values {
		if strings.Contains(result, value) {
			return true, nil
		}
	}

	return false, nil
}

func validateFilterRegex(values []string, result string) (bool, error) {
	for _, value := range values {
		r, err := regexp.Compile(value)
		if err != nil {
			return false, fmt.Errorf("failed to compile regex: %s", err)
		}

		if r.MatchString(result) {
			return true, nil
		}
	}

	return false, nil
}
