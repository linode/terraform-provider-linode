package helper

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SDKv2UnwrapOptionalAttr returns a pointer to a value from the given resource
// data. If the value is known the function returns a pointer to it, else it returns
// nil.
func SDKv2UnwrapOptionalAttr[T any](d *schema.ResourceData, path string) *T {
	value, ok := d.GetOk(path)
	if !ok || value == nil {
		return nil
	}

	castedValue := value.(T)

	return &castedValue
}
