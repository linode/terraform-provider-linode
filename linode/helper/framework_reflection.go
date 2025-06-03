package helper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FrameworkModelToObjectType[T any](ctx context.Context) (types.ObjectType, error) {
	reflectedType := reflect.TypeFor[T]()

	// Deref pointers if necessary
	for reflectedType.Kind() == reflect.Ptr {
		reflectedType = reflectedType.Elem()
	}

	if reflectedType.Kind() != reflect.Struct {
		return types.ObjectType{}, fmt.Errorf("expected a struct, got %s", reflectedType.Kind().String())
	}

	reflectedAttrValue := reflect.TypeFor[attr.Value]()

	resultAttributes := make(map[string]attr.Type)

	for fieldIndex := range reflectedType.NumField() {
		reflectedField := reflectedType.Field(fieldIndex)

		tfsdkTag, ok := reflectedField.Tag.Lookup("tfsdk")
		if !ok {
			continue
		}

		if _, ok := resultAttributes[tfsdkTag]; ok {
			return types.ObjectType{}, fmt.Errorf("found duplicate tfsdk tag: %s", tfsdkTag)
		}

		// Deref pointers if necessary
		reflectedFieldType := reflectedField.Type

		for reflectedFieldType.Kind() == reflect.Ptr {
			reflectedFieldType = reflectedType.Elem()
		}
		if !reflectedFieldType.Implements(reflectedAttrValue) {
			return types.ObjectType{}, fmt.Errorf(
				"field %s does not implement attr.Value: %s", reflectedFieldType.Name(), reflectedFieldType.String(),
			)
		}

		attrType := reflect.New(reflectedFieldType).Interface().(attr.Value).Type(ctx)

		resultAttributes[tfsdkTag] = attrType
	}

	return types.ObjectType{
		AttrTypes: resultAttributes,
	}, nil
}
