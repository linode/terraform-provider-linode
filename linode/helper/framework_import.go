package helper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportStatePassthroughInt64ID allows for the automatic importing of resources
// through an int64 ID attribute. This is necessary as many Linode resources
// use integers rather than strings as unique identifiers.
func ImportStatePassthroughInt64ID(
	ctx context.Context,
	attrPath path.Path,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportStatePassthroughInt64ID path must be set to a valid attribute path.",
		)
	}

	intID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Integer Value for Import ID",
			fmt.Sprintf("\"%s\" is not a valid integer value", req.ID),
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, intID)...)
}

type IDTypeConverter func(values ...string) ([]any, diag.Diagnostics)

func AllInt64(values ...string) ([]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := make([]any, len(values))
	for i, v := range values {
		int64Val := StringToInt64(v, &diags)
		result[i] = int64Val
	}
	return result, nil
}

func AllString(values ...string) ([]any, diag.Diagnostics) {
	result := make([]any, len(values))
	for i, v := range values {
		result[i] = v
	}
	return result, nil
}
