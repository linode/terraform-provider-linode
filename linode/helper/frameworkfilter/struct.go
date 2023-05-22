package frameworkfilter

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// resolveStructFieldByJSON resolves the struct field resolves the StructField
// with the given JSON tag.
func resolveStructFieldByJSON(val any, field string) (reflect.StructField, diag.Diagnostic) {
	rType := reflect.TypeOf(val)

	for i := 0; i < rType.NumField(); i++ {
		currentField := rType.Field(i)
		if tag, ok := currentField.Tag.Lookup("json"); ok && tag == field {
			return currentField, nil
		}
	}

	return reflect.StructField{}, diag.NewErrorDiagnostic(
		"Missing field",
		fmt.Sprintf("Failed to find field %s in struct.", field),
	)
}

// resolveStructValueByJSON resolves the corresponding value of a struct field
// given a JSON tag.
func resolveStructValueByJSON(val any, field string) (any, diag.Diagnostic) {
	structField, d := resolveStructFieldByJSON(val, field)
	if d != nil {
		return nil, d
	}

	targetField := reflect.ValueOf(val).FieldByName(structField.Name)

	if !targetField.IsValid() {
		return nil, diag.NewErrorDiagnostic(
			"Field not found",
			fmt.Sprintf("Could not find JSON tag in target struct: %s", field),
		)
	}

	return targetField.Interface(), nil
}
