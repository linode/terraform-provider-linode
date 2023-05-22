package frameworkfilter

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"reflect"
	"time"
)

// GetLatestCreated is a helper function that returns the latest
// create entry in the input slice.
func (f Config) GetLatestCreated(elems []any, field string) (any, diag.Diagnostic) {
	if len(elems) < 1 {
		return nil, nil
	}

	timeField, d := resolveStructFieldByJSON(elems[0], field)
	if d != nil {
		return nil, d
	}

	newestElem := elems[0]

	for _, elem := range elems {
		newestElemCreated, d := getCreatedTime(newestElem, timeField.Name)
		if d != nil {
			return nil, d
		}

		currentElemCreated, d := getCreatedTime(elem, timeField.Name)
		if d != nil {
			return nil, d
		}

		if currentElemCreated.After(newestElemCreated) {
			newestElem = elem
		}
	}

	return newestElem, nil
}

// getCreatedTime parses a time value from the given elem and field.
func getCreatedTime(elem any, attr string) (time.Time, diag.Diagnostic) {
	val := reflect.ValueOf(elem).FieldByName(attr)

	// Deref any pointers
	for val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	result, ok := val.Interface().(time.Time)
	if !ok {
		return time.Time{}, diag.NewErrorDiagnostic(
			"Field has incorrect type",
			fmt.Sprintf("Field %s is not type time.Time", attr),
		)
	}

	return result, nil
}
