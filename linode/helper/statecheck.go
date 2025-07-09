package helper

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/compare"
)

var _ compare.ValueComparer = typeAgnosticComparer{}

type typeAgnosticComparer struct{}

func (typeAgnosticComparer) CompareValues(values ...any) error {
	var lastNormalizedValue *string

	for _, value := range values {
		normalizedValue := ""

		switch value.(type) {
		case string:
			normalizedValue = value.(string)
		case json.Number:
			normalizedValue = value.(json.Number).String()
		default:
			return fmt.Errorf("unsupported type for comparisons: %T", value)
		}

		if lastNormalizedValue != nil && normalizedValue != *lastNormalizedValue {
			return fmt.Errorf("found difference: %s != %s", normalizedValue, *lastNormalizedValue)
		}

		lastNormalizedValue = &normalizedValue
	}

	return nil
}

func TypeAgnosticComparer() typeAgnosticComparer {
	return typeAgnosticComparer{}
}
