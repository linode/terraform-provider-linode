package helper

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDKv2UnwrapOptionalConfigAttr returns a pointer to a value from the given path
// only if the value was explicitly defined.
//
// If the path cannot be resolved, the function will return nil as
// this may be valid case for nested fields.
//
// This is useful to simplifying the logic around checking whether a nested field has been
// explicitly defined by a user.
func SDKv2UnwrapOptionalConfigAttr[T any](ctx context.Context, d *schema.ResourceData, path string) *T {
	ctyPath := SDKv2PathToCtyPath(path)

	ctyValue, err := ctyPath.Apply(d.GetRawConfig())
	if err != nil {
		// This is a valid case, but we should log it just in case something was expected to be
		// resolved
		tflog.Trace(ctx, "failed to resolve path, assuming undefined", map[string]any{
			"path": path,
		})
		return nil
	}

	// If the value hasn't been defined by the user, return nil
	if !ctyValue.IsKnown() || ctyValue.IsNull() {
		return nil
	}

	// If the value was defined explicitly as nil, return nil
	value := d.Get(path)
	if value == nil {
		return nil
	}

	castedValue := value.(T)
	return &castedValue
}

// SDKv2PathToCtyPath converts an SDKv2-style path (e.g. foo.0.bar)
// to a cty-style path object.
func SDKv2PathToCtyPath(path string) cty.Path {
	segments := strings.Split(path, ".")
	result := make(cty.Path, len(segments))

	for i, segment := range segments {
		// If this is a number, treat it as an index
		if value, err := strconv.ParseInt(segment, 10, 64); err == nil {
			result[i] = cty.IndexStep{
				Key: cty.NumberIntVal(value),
			}
			continue
		}

		// Otherwise, treat is as an attribute
		result[i] = cty.GetAttrStep{
			Name: segment,
		}
	}

	return result
}
