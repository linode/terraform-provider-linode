package frameworkfilter

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const EXACT = "exact"
const SUBSTRING = "substring"
const REGEX = "regex"

// applyLocalFiltering handles filtering for fields that are not
// API-filterable.
func (f Config) applyLocalFiltering(
	filterSet []FilterModel, data []any,
) ([]any, diag.Diagnostic) {
	result := make([]any, 0)

	for _, elem := range data {
		match, d := f.matchesFilter(filterSet, elem)
		if d != nil {
			return nil, d
		}

		// This element was filtered out
		if !match {
			continue
		}

		result = append(result, elem)
	}

	return result, nil
}

// matchesFilter checks whether an object matches the given filter set.
func (f Config) matchesFilter(
	filterSet []FilterModel,
	elem any,
) (bool, diag.Diagnostic) {
	for _, filter := range filterSet {
		filterName := filter.Name.ValueString()

		// Skip if this field should be filtered at an API level
		if f[filterName].APIFilterable && isExactMatchFilter(filter) {
			continue
		}

		// Grab the field from the input struct
		matchingField, d := resolveStructValueByJSON(elem, filterName)
		if d != nil {
			return false, d
		}

		// Check whether the field matches the filter
		match, d := f.checkFieldMatchesFilter(matchingField, filter)
		if d != nil {
			return false, d
		}

		// No match for this filter; return
		if !match {
			return false, nil
		}
	}

	return true, nil
}

// checkFieldMatchesFilter checks whether an individual field
// meets the condition for the given filter.
func (f Config) checkFieldMatchesFilter(
	field any,
	filter FilterModel,
) (bool, diag.Diagnostic) {
	rField := reflect.ValueOf(field)

	// Recursively filter on list elements (tags, capabilities, etc.)
	if rField.Kind() == reflect.Slice {
		for i := 0; i < rField.Len(); i++ {
			match, d := f.checkFieldMatchesFilter(rField.Index(i).Interface(), filter)
			if d != nil {
				return false, d
			}

			if match {
				return true, nil
			}
		}

		return false, nil
	}

	normalizedValue, d := normalizeValue(field)
	if d != nil {
		return false, d
	}

	// Run the corresponding validation method
	result := false
	d = nil

	switch strings.ToLower(filter.MatchBy.ValueString()) {
	case EXACT, "":
		result = checkFilterExact(filter.Values, normalizedValue)
	case SUBSTRING, "sub":
		result, d = checkFilterSubString(filter.Values, normalizedValue)
	case REGEX, "re":
		result, d = checkFilterRegex(filter.Values, normalizedValue)
	}

	return result, d
}

// normalizeValue converts the given field into a comparable string.
func normalizeValue(field any) (string, diag.Diagnostic) {
	rField := reflect.ValueOf(field)

	// Dereference if the value is a pointer
	for rField.Kind() == reflect.Pointer {
		// Null pointer; assume empty
		if rField.IsNil() {
			return "", nil
		}

		rField = reflect.Indirect(rField)
	}
	switch rField.Kind() {
	case reflect.String:
		return rField.String(), nil
	case reflect.Int, reflect.Int64:
		return strconv.FormatInt(rField.Int(), 10), nil
	case reflect.Bool:
		return strconv.FormatBool(rField.Bool()), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rField.Float(), 'f', 0, 64), nil
	default:
		return "", diag.NewErrorDiagnostic(
			"Invalid field type",
			fmt.Sprintf("Invalid type for field: %s", rField.Type().String()),
		)
	}
}

func checkFilterExact(values []types.String, actualValue string) bool {
	for _, value := range values {
		if reflect.DeepEqual(actualValue, value.ValueString()) {
			return true
		}
	}

	return false
}

func checkFilterSubString(values []types.String, actualValue string) (bool, diag.Diagnostic) {
	for _, value := range values {
		if strings.Contains(actualValue, value.ValueString()) {
			return true, nil
		}
	}

	return false, nil
}

func checkFilterRegex(values []types.String, actualValue string) (bool, diag.Diagnostic) {
	for _, value := range values {
		r, err := regexp.Compile(value.ValueString())
		if err != nil {
			return false, diag.NewErrorDiagnostic(
				"failed to compile regex",
				err.Error(),
			)
		}

		if r.MatchString(actualValue) {
			return true, nil
		}
	}

	return false, nil
}

func isExactMatchFilter(filter FilterModel) bool {
	return filter.MatchBy.ValueString() == EXACT || filter.MatchBy.IsNull()
}
