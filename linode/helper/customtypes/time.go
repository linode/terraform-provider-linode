package customtypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringValuableWithSemanticEquals = RFC3339TimeStringValue{}
var _ basetypes.StringTypable = RFC3339TimeStringType{}
var _ xattr.TypeWithValidate = RFC3339TimeStringType{}

type RFC3339TimeStringValue struct {
	basetypes.StringValue
}

func (v RFC3339TimeStringValue) Equal(o attr.Value) bool {
	other, ok := o.(RFC3339TimeStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v RFC3339TimeStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(RFC3339TimeStringValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return helper.CompareRFC3339TimeStrings(
		v.StringValue.ValueString(),
		newValue.ValueString(),
	), diags
}

func (v RFC3339TimeStringValue) Type(ctx context.Context) attr.Type {
	return RFC3339TimeStringType{}
}

type RFC3339TimeStringType struct {
	basetypes.StringType
}

func (t RFC3339TimeStringType) String() string {
	return "RFC3339TimeStringType"
}

func (t RFC3339TimeStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// CustomStringValue defined in the value type section
	value := RFC3339TimeStringValue{
		StringValue: in,
	}
	return value, nil
}

func (t RFC3339TimeStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t RFC3339TimeStringType) ValueType(ctx context.Context) attr.Value {
	return RFC3339TimeStringValue{}
}

func (t RFC3339TimeStringType) Equal(o attr.Type) bool {
	other, ok := o.(RFC3339TimeStringType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t RFC3339TimeStringType) Validate(ctx context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
	if value.IsNull() || !value.IsKnown() {
		return nil
	}

	var diags diag.Diagnostics
	var valueString string

	if err := value.As(&valueString); err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid Terraform Value",
			"An unexpected error occurred while attempting to convert a Terraform value to a string. "+
				"This generally is an issue with the provider schema implementation. "+
				"Please contact the provider developers.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	if _, err := time.Parse(time.RFC3339, valueString); err != nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid RFC 3339 String Value",
			"An unexpected error occurred while converting a string value that was expected to be RFC 3339 format. "+
				"The RFC 3339 string format is YYYY-MM-DDTHH:MM:SSZ, such as 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00.\n\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString+"\n"+
				"Error: "+err.Error(),
		)

		return diags
	}

	return diags
}
