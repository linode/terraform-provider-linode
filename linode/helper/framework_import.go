package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// IDTypeConverter represents a function that converts a given value into
// a given type
type IDTypeConverter func(value string) (any, diag.Diagnostics)

// ImportableID represents a single ID attribute importable
type ImportableID struct {
	TypeConverter IDTypeConverter
	Name          string
}

func IDTypeConverterString(value string) (any, diag.Diagnostics) {
	return value, nil
}

func IDTypeConverterInt64(value string) (any, diag.Diagnostics) {
	var d diag.Diagnostics
	result := StringToInt64(value, &d)
	return result, d
}

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

// ImportStateWithMultipleIDs allows framework resources with multiple IDs
// (e.g. `child_id, parent_id`) be imported.
func ImportStateWithMultipleIDs(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
	idFields []ImportableID,
) {
	// Compute a readable list of required IDs
	idFieldNames := make([]string, len(idFields))
	for i, field := range idFields {
		idFieldNames[i] = field.Name
	}

	unexpectedIDsErrorMsg := fmt.Sprintf(
		"Expected import identifier with format: %s. Got: %q",
		strings.Join(idFieldNames, ", "), req.ID,
	)

	// Make sure we support spaces in the ID just in case :)
	fullID := strings.ReplaceAll(req.ID, " ", "")

	idParts := strings.Split(fullID, ",")

	if len(idParts) != len(idFields) {
		resp.Diagnostics.AddError("Unexpected Import Identifier", unexpectedIDsErrorMsg)
		return
	}

	for i, id := range idParts {
		if id == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier", unexpectedIDsErrorMsg,
			)
			return
		}

		field := idFields[i]

		valueConverted, d := field.TypeConverter(id)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(
			resp.State.SetAttribute(ctx, path.Root(field.Name), valueConverted)...,
		)
	}
}
