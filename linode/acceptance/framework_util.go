package acceptance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// FrameworkUnwrap runs the given framework function, triggering a test failure
// if the returned diagnostics have an error.
func FrameworkUnwrap(t *testing.T, inner func() diag.Diagnostics) {
	t.Helper()

	if d := inner(); d.HasError() {
		t.Fatal(d.Errors())
	}
}

// FrameworkObjectAs expands the given types.Object to the given return type,
// causing the test to fail if an error occurs.
func FrameworkObjectAs[O any](t *testing.T, object types.Object) O {
	t.Helper()

	var result O

	FrameworkUnwrap(t, func() diag.Diagnostics {
		return object.As(
			context.Background(),
			&result,
			basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			},
		)
	})

	return result
}
