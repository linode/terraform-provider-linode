package customtypes

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringTypable                    = IPAddrStringType{}
	_ basetypes.StringValuableWithSemanticEquals = IPAddrStringValue{}
	_ xattr.TypeWithValidate                     = IPAddrStringType{}
)

// IPAddrStringType represents the type of an attribute that
// can be either an IPv4 address or an IPv6 address.
// Additionally, this type implements validation the responding validation.
type IPAddrStringType struct {
	basetypes.StringType
}

func (t IPAddrStringType) Equal(o attr.Type) bool {
	other, ok := o.(IPAddrStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t IPAddrStringType) String() string {
	return "IPAddrStringType"
}

func (t IPAddrStringType) ValueFromString(
	ctx context.Context,
	in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPAddrStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t IPAddrStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t IPAddrStringType) ValueType(ctx context.Context) attr.Value {
	return IPAddrStringValue{}
}

func (t IPAddrStringType) Validate(ctx context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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

	ip := net.ParseIP(valueString)

	if ip == nil {
		diags.AddAttributeError(
			valuePath,
			"Invalid IPv4/IPv6 Address Value",
			"An unexpected error occurred while parsing a string value as an IPv4 or IPv6 address.\n"+
				"Path: "+valuePath.String()+"\n"+
				"Given Value: "+valueString+"\n",
		)

		return diags
	}

	return diags
}

var _ basetypes.StringValuable = IPAddrStringValue{}

// IPAddrStringValue represents a string that can be an IPv4 or IPv6 address.
// This value implements semantic equality checks for semantically equal IPv6 addresses.
type IPAddrStringValue struct {
	basetypes.StringValue
}

func (v IPAddrStringValue) Equal(o attr.Value) bool {
	other, ok := o.(IPAddrStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v IPAddrStringValue) Type(ctx context.Context) attr.Type {
	return IPAddrStringType{}
}

func (v IPAddrStringValue) StringSemanticEquals(
	ctx context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(IPAddrStringValue)
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

	return net.ParseIP(v.ValueString()).Equal(net.ParseIP(newValue.ValueString())), nil
}

func IPAddrValue(value string) IPAddrStringValue {
	return IPAddrStringValue{
		StringValue: types.StringValue(value),
	}
}
