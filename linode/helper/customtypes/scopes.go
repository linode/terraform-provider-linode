package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringTypable                    = LinodeScopesStringType{}
	_ basetypes.StringValuableWithSemanticEquals = LinodeScopesStringValue{}
)

type LinodeScopesStringType struct {
	basetypes.StringType
}

func (t LinodeScopesStringType) Equal(o attr.Type) bool {
	other, ok := o.(LinodeScopesStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t LinodeScopesStringType) String() string {
	return "LinodeScopesStringType"
}

func (t LinodeScopesStringType) ValueFromString(
	ctx context.Context,
	in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	// LinodeScopesStringValue defined in the value type section
	value := LinodeScopesStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t LinodeScopesStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t LinodeScopesStringType) ValueType(ctx context.Context) attr.Value {
	return LinodeScopesStringValue{}
}

var _ basetypes.StringValuable = LinodeScopesStringValue{}

type LinodeScopesStringValue struct {
	basetypes.StringValue
}

func (v LinodeScopesStringValue) Equal(o attr.Value) bool {
	other, ok := o.(LinodeScopesStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v LinodeScopesStringValue) Type(ctx context.Context) attr.Type {
	return LinodeScopesStringType{}
}

// LinodeScopesStringValue defined in the value type section
// Ensure the implementation satisfies the expected interfaces

func (v LinodeScopesStringValue) StringSemanticEquals(
	ctx context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(LinodeScopesStringValue)
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

	return helper.CompareScopes(v.StringValue.ValueString(), newValue.ValueString()), diags
}
