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

// Assign StringNull() safely without throwing error. e.g. new value: .rev_note: was null, but now cty.StringVal("")
func getValueIfNotNull(val string) basetypes.StringValue {
	res := types.StringValue(val)

	if res == types.StringValue("") {
		res = types.StringNull()
	}

	return res
}

func (data *StackScriptModel) parseStackScript(
	ctx context.Context,
	stackscript *linodego.Stackscript,
) diag.Diagnostics {
	data.ID = types.StringValue(strconv.Itoa(stackscript.ID))
	data.Label = types.StringValue(stackscript.Label)
	data.Script = types.StringValue(stackscript.Script)
	data.Description = types.StringValue(stackscript.Description)
	data.RevNote = getValueIfNotNull(stackscript.RevNote)
	data.IsPublic = types.BoolValue(stackscript.IsPublic)

	images, err := types.ListValueFrom(ctx, types.StringType, stackscript.Images)
	if err != nil {
		return err
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
		if err != nil {
			return err
		}

		data.UserDefinedFields = *udf
	}

	return nil
}

// flattenUserDefinedFields flattens a list of linodego UDF objects into a basetypes.ListValue
func flattenUserDefinedFields(udf []linodego.StackscriptUDF) (*basetypes.ListValue, diag.Diagnostics) {
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
		if err != nil {
			return nil, err
		}

		resultList[i] = obj
	}

	result, err := basetypes.NewListValue(
		udfObjectType,
		resultList,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
