package customdiffs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ComputedWithDefault allows a computed field to have an explicit default if
// it is not defined by the user.
//
// This is hacky but allows us to avoid a breaking change on fields using
// CaseInsensitiveSet (computed) when not specifying a field.
func ComputedWithDefault[T any](field string, defaultValue T) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
		if !diff.GetRawConfig().GetAttr(field).IsNull() {
			return nil
		}

		return diff.SetNew(field, defaultValue)
	}
}
