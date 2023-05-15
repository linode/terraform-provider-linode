package stackscript

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// StackScriptModel describes the Terraform resource data model to match the
// resource schema.
type StackScriptModel struct {
	ID                types.String        `tfsdk:"id"`
	Label             types.String        `tfsdk:"label"`
	Script            types.String        `tfsdk:"script"`
	Description       types.String        `tfsdk:"description"`
	RevNote           types.String        `tfsdk:"rev_note"`
	IsPublic          types.Bool          `tfsdk:"is_public"`
	Images            types.List          `tfsdk:"images"`
	DeploymentsActive types.Int64         `tfsdk:"deployments_active"`
	UserGravatarID    types.String        `tfsdk:"user_gravatar_id"`
	DeploymentsTotal  types.Int64         `tfsdk:"deployments_total"`
	Username          types.String        `tfsdk:"username"`
	Created           types.String        `tfsdk:"created"`
	Updated           types.String        `tfsdk:"updated"`
	UserDefinedFields basetypes.ListValue `tfsdk:"user_defined_fields"`
}

func (data *StackScriptModel) parseNonComputedAttributes(
	ctx context.Context,
	stackscript *linodego.Stackscript,
) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	data.Label = types.StringValue(stackscript.Label)
	data.Script = types.StringValue(stackscript.Script)
	data.Description = types.StringValue(stackscript.Description)

	// These fields have schema-defined defaults,
	// so this should not cause comparison issues
	data.RevNote = types.StringValue(stackscript.RevNote)
	data.IsPublic = types.BoolValue(stackscript.IsPublic)

	// Only update the images return if there is a change
	var plannedImages []string
	diagnostics.Append(data.Images.ElementsAs(ctx, &plannedImages, false)...)
	if diagnostics.HasError() {
		return diagnostics
	}

	if !helper.StringListElementsEqual(plannedImages, stackscript.Images) {
		remoteImages, err := types.ListValueFrom(ctx, types.StringType, stackscript.Images)
		diagnostics.Append(err...)
		if diagnostics.HasError() {
			return diagnostics
		}

		data.Images = remoteImages
	}

	return diagnostics
}

func (data *StackScriptModel) parseComputedAttributes(
	ctx context.Context,
	stackscript *linodego.Stackscript,
) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(stackscript.ID))

	images, err := types.ListValueFrom(ctx, types.StringType, stackscript.Images)
	diagnostics.Append(err...)
	if diagnostics.HasError() {
		return diagnostics
	}

	data.Images = images
	data.DeploymentsActive = types.Int64Value(int64(stackscript.DeploymentsActive))
	data.UserGravatarID = types.StringValue(stackscript.UserGravatarID)
	data.DeploymentsTotal = types.Int64Value(int64(stackscript.DeploymentsTotal))
	data.Username = types.StringValue(stackscript.Username)
	data.Created = types.StringValue(stackscript.Created.Format(time.RFC3339))
	data.Updated = types.StringValue(stackscript.Updated.Format(time.RFC3339))

	if stackscript.UserDefinedFields != nil {
		udf, err := flattenUserDefinedFields(*stackscript.UserDefinedFields)
		diagnostics.Append(err...)
		if diagnostics.HasError() {
			return diagnostics
		}

		data.UserDefinedFields = *udf
	}

	return diagnostics
}

// flattenUserDefinedFields flattens a list of linodego UDF objects into a basetypes.ListValue
func flattenUserDefinedFields(udf []linodego.StackscriptUDF) (*basetypes.ListValue, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	resultList := make([]attr.Value, len(udf))

	for i, field := range udf {
		valueMap := make(map[string]attr.Value)
		valueMap["label"] = types.StringValue(field.Label)
		valueMap["name"] = types.StringValue(field.Name)
		valueMap["example"] = types.StringValue(field.Example)
		valueMap["one_of"] = types.StringValue(field.OneOf)
		valueMap["many_of"] = types.StringValue(field.ManyOf)
		valueMap["default"] = types.StringValue(field.Default)

		obj, err := types.ObjectValue(udfObjectType.AttrTypes, valueMap)
		diagnostics.Append(err...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		resultList[i] = obj
	}

	result, err := basetypes.NewListValue(
		udfObjectType,
		resultList,
	)
	diagnostics.Append(err...)
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return &result, nil
}
