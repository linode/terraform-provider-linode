package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FrameworkAttemptRemoveResourceForEmptyID implements a
// temporary workaround for a Crossplane/Upjet issue (TPT-2408).
// Returns true if the resource was removed from state, else false.
func FrameworkAttemptRemoveResourceForEmptyID(
	ctx context.Context,
	id types.String,
	resp *resource.ReadResponse,
) bool {
	if id.ValueString() != "" {
		return false
	}

	resp.Diagnostics.AddWarning(
		"Removing Resource From State",
		"This resource is being implicitly removed from the Terraform state because "+
			"its ID attribute is empty.",
	)
	resp.State.RemoveResource(ctx)

	return true
}
