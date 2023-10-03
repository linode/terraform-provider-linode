package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SetLogFieldBulk allows for setting multiple logger fields at a time.
func SetLogFieldBulk(ctx context.Context, fields map[string]any) context.Context {
	result := ctx

	for k, v := range fields {
		result = tflog.SetField(result, k, v)
	}

	return result
}
