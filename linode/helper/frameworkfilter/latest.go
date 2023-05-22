package frameworkfilter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// GetLatestCreated is a helper function that returns the latest
// create entry in the input slice.
func (f Config) GetLatestCreated(elems []any, structField string) (any, diag.Diagnostic) {
	if len(elems) < 1 {
		return nil, nil
	}

	newestElem := elems[0]
	newestCreated, d := getCreatedTime(newestElem, structField)
	if d != nil {
		return nil, d
	}

	for _, elem := range elems {
		currentElemCreated, d := getCreatedTime(elem, structField)
		if d != nil {
			return nil, d
		}

		if currentElemCreated.After(newestCreated) {
			newestElem = elem
			newestCreated = currentElemCreated
		}
	}

	return newestElem, nil
}

// GetLatestVersion gets the latest version of the given struct
func (f Config) GetLatestVersion(elems []any, structField string) (any, diag.Diagnostic) {
	if len(elems) < 1 {
		return nil, nil
	}

	newestElem := elems[0]
	newestVersion, d := getVersion(newestElem, structField)
	if d != nil {
		return nil, d
	}

	for _, elem := range elems {
		currentVersion, d := getVersion(elem, structField)
		if d != nil {
			return nil, d
		}

		if currentVersion.GreaterThan(newestVersion) {
			newestElem = elem
			newestVersion = currentVersion
		}
	}

	return newestElem, nil
}

// getCreatedTime parses a time value from the given elem and field.
func getCreatedTime(elem any, attr string) (time.Time, diag.Diagnostic) {
	val := reflect.ValueOf(elem).FieldByName(attr)
	if !val.IsValid() {
		return time.Time{}, diag.NewErrorDiagnostic(
			"Field does not exist",
			fmt.Sprintf("Field %s does not exist in struct", attr),
		)
	}

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

// getVersion parses a version value from the given elem and field.
func getVersion(elem any, attr string) (*version.Version, diag.Diagnostic) {
	val := reflect.ValueOf(elem).FieldByName(attr)
	if !val.IsValid() {
		return nil, diag.NewErrorDiagnostic(
			"Field does not exist",
			fmt.Sprintf("Field %s does not exist in struct", attr),
		)
	}

	// Deref any pointers
	for val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// Check the type
	if val.Kind() != reflect.String {
		return nil, diag.NewErrorDiagnostic(
			"Field has incorrect type",
			fmt.Sprintf("Field %s is not type string", attr),
		)
	}

	// Parse the version
	result, err := version.NewVersion(val.String())
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			"Failed to parse version",
			err.Error(),
		)
	}

	return result, nil
}
