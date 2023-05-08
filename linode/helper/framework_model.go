package helper

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// ShouldModelUpdate is a helper function that checks whether a model has
// been updated. This is useful for simplifying resource update logic.
//
// NOTE: Only fields marked with the `linode_mutable:"true"` tag will be compared.
func ShouldModelUpdate[T any](model1, model2 T) (bool, error) {
	reflectModel1 := reflect.ValueOf(model1)
	reflectModel2 := reflect.ValueOf(model2)

	// Deref the pointers if necessary
	if reflectModel1.Kind() == reflect.Ptr {
		reflectModel1 = reflect.Indirect(reflectModel1)
		reflectModel2 = reflect.Indirect(reflectModel2)
	}

	modelType := reflectModel1.Type()

	for i := 0; i < modelType.NumField(); i++ {
		currentField := modelType.Field(i)

		// Check that the field is mutable
		mutableTag, ok := currentField.Tag.Lookup("linode_mutable")
		if !ok {
			continue
		}

		mutable, err := strconv.ParseBool(mutableTag)
		if err != nil {
			return false, fmt.Errorf("failed to parse linode_mutable tag: %w", err)
		}

		if !mutable {
			continue
		}

		// Get the name of the field
		fieldName := currentField.Name

		// Ensure the field implements attr.Value
		model1Field := reflectModel1.FieldByName(fieldName)
		if !model1Field.IsValid() {
			return false, fmt.Errorf("the model type is missing the specified %s field", fieldName)
		}

		if !model1Field.Type().Implements(reflect.TypeOf((*attr.Value)(nil)).Elem()) {
			return false, fmt.Errorf("field %s does not implement attr.Value", fieldName)
		}

		model2Field := reflectModel2.FieldByName(fieldName)

		// Cast and compare
		model1FieldValue := model1Field.Interface().(attr.Value)
		model2FieldValue := model2Field.Interface().(attr.Value)

		if !model1FieldValue.Equal(model2FieldValue) {
			return true, nil
		}
	}

	return false, nil
}
