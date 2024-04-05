package stackscript

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// StackScriptModel describes the older Terraform resource data model to match the
// resource schema of Terraform Provider for Linode v1.
// The only difference is the format of time for created and updated attributes
type StackScriptModelV0 struct {
	ID                types.String `tfsdk:"id"`
	Label             types.String `tfsdk:"label"`
	Script            types.String `tfsdk:"script"`
	Description       types.String `tfsdk:"description"`
	RevNote           types.String `tfsdk:"rev_note"`
	IsPublic          types.Bool   `tfsdk:"is_public"`
	Images            types.Set    `tfsdk:"images"`
	DeploymentsActive types.Int64  `tfsdk:"deployments_active"`
	UserGravatarID    types.String `tfsdk:"user_gravatar_id"`
	DeploymentsTotal  types.Int64  `tfsdk:"deployments_total"`
	Username          types.String `tfsdk:"username"`
	Created           types.String `tfsdk:"created"`
	Updated           types.String `tfsdk:"updated"`
	UserDefinedFields types.List   `tfsdk:"user_defined_fields"`
}

// StackScriptModel describes the Terraform resource data model to match the
// resource schema.
type StackScriptModel struct {
	ID                types.String      `tfsdk:"id"`
	Label             types.String      `tfsdk:"label"`
	Script            types.String      `tfsdk:"script"`
	Description       types.String      `tfsdk:"description"`
	RevNote           types.String      `tfsdk:"rev_note"`
	IsPublic          types.Bool        `tfsdk:"is_public"`
	Images            types.Set         `tfsdk:"images"`
	DeploymentsActive types.Int64       `tfsdk:"deployments_active"`
	UserGravatarID    types.String      `tfsdk:"user_gravatar_id"`
	DeploymentsTotal  types.Int64       `tfsdk:"deployments_total"`
	Username          types.String      `tfsdk:"username"`
	Created           timetypes.RFC3339 `tfsdk:"created"`
	Updated           timetypes.RFC3339 `tfsdk:"updated"`
	UserDefinedFields types.List        `tfsdk:"user_defined_fields"`
}

func (data *StackScriptModel) FlattenStackScript(
	stackscript *linodego.Stackscript,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(stackscript.ID), preserveKnown)
	data.Label = helper.KeepOrUpdateString(data.Label, stackscript.Label, preserveKnown)
	data.Script = helper.KeepOrUpdateString(data.Script, stackscript.Script, preserveKnown)
	data.Description = helper.KeepOrUpdateString(
		data.Description, stackscript.Description, preserveKnown,
	)
	data.RevNote = helper.KeepOrUpdateString(data.RevNote, stackscript.RevNote, preserveKnown)
	data.IsPublic = helper.KeepOrUpdateBool(data.IsPublic, stackscript.IsPublic, preserveKnown)
	data.DeploymentsActive = helper.KeepOrUpdateInt64(
		data.DeploymentsActive, int64(stackscript.DeploymentsActive), preserveKnown,
	)
	data.UserGravatarID = helper.KeepOrUpdateString(
		data.UserGravatarID, stackscript.UserGravatarID, preserveKnown,
	)
	data.DeploymentsTotal = helper.KeepOrUpdateInt64(
		data.DeploymentsTotal, int64(stackscript.DeploymentsTotal), preserveKnown,
	)
	data.Username = helper.KeepOrUpdateString(data.Username, stackscript.Username, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(
		data.Created, timetypes.NewRFC3339TimePointerValue(stackscript.Created), preserveKnown,
	)
	data.Updated = helper.KeepOrUpdateValue(
		data.Updated, timetypes.NewRFC3339TimePointerValue(stackscript.Updated), preserveKnown,
	)

	data.Images = helper.KeepOrUpdateStringSet(
		data.Images, stackscript.Images, preserveKnown, &diags,
	)

	if stackscript.UserDefinedFields == nil {
		data.UserDefinedFields = helper.KeepOrUpdateValue(
			data.UserDefinedFields, types.ListNull(udfObjectType), preserveKnown,
		)
	} else {
		udf, err := flattenUserDefinedFields(*stackscript.UserDefinedFields)
		diags.Append(err...)
		if diags.HasError() {
			return diags
		}

		if udf != nil {
			data.UserDefinedFields = helper.KeepOrUpdateValue(
				data.UserDefinedFields, *udf, preserveKnown,
			)
		}
	}

	return diags
}

func (data *StackScriptModel) CopyFrom(other StackScriptModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Label = helper.KeepOrUpdateValue(data.Label, other.Label, preserveKnown)
	data.Script = helper.KeepOrUpdateValue(data.Script, other.Script, preserveKnown)
	data.Description = helper.KeepOrUpdateValue(data.Description, other.Description, preserveKnown)
	data.RevNote = helper.KeepOrUpdateValue(data.RevNote, other.RevNote, preserveKnown)
	data.IsPublic = helper.KeepOrUpdateValue(data.IsPublic, other.IsPublic, preserveKnown)
	data.Images = helper.KeepOrUpdateValue(data.Images, other.Images, preserveKnown)
	data.DeploymentsActive = helper.KeepOrUpdateValue(data.DeploymentsActive, other.DeploymentsActive, preserveKnown)
	data.UserGravatarID = helper.KeepOrUpdateValue(data.UserGravatarID, other.UserGravatarID, preserveKnown)
	data.DeploymentsTotal = helper.KeepOrUpdateValue(data.DeploymentsTotal, other.DeploymentsTotal, preserveKnown)
	data.Username = helper.KeepOrUpdateValue(data.Username, other.Username, preserveKnown)
	data.Created = helper.KeepOrUpdateValue(data.Created, other.Created, preserveKnown)
	data.Updated = helper.KeepOrUpdateValue(data.Updated, other.Updated, preserveKnown)
	data.UserDefinedFields = helper.KeepOrUpdateValue(data.UserDefinedFields, other.UserDefinedFields, preserveKnown)
}

// flattenUserDefinedFields flattens a list of linodego UDF objects into a types.List
func flattenUserDefinedFields(udf []linodego.StackscriptUDF) (*types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if udf == nil {
		nullList := types.ListNull(udfObjectType)
		return &nullList, diags
	}

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
		diags.Append(err...)
		if diags.HasError() {
			return nil, diags
		}

		resultList[i] = obj
	}

	result, err := types.ListValue(
		udfObjectType,
		resultList,
	)
	diags.Append(err...)
	if diags.HasError() {
		return nil, diags
	}

	return &result, nil
}
