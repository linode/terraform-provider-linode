package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	// Ensure the implementation satisfies the expected interfaces
	_ basetypes.StringTypable                    = LinodeAutoAllocIPType{}
	_ basetypes.StringValuableWithSemanticEquals = LinodeAutoAllocIPValue{}
)

// LinodeAutoAllocIPType represents the type of an attribute that accepts either a
// specific IP address or the literal "auto" to allocate one automatically.
//
// It implements semantic equality so that a configured value of "auto" does not
// cause a perpetual diff against the concrete IP address returned by the API.
type LinodeAutoAllocIPType struct {
	basetypes.StringType
}

func (t LinodeAutoAllocIPType) Equal(o attr.Type) bool {
	other, ok := o.(LinodeAutoAllocIPType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t LinodeAutoAllocIPType) String() string {
	return "LinodeAutoAllocIPType"
}

func (t LinodeAutoAllocIPType) ValueFromString(
	ctx context.Context,
	in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := LinodeAutoAllocIPValue{
		StringValue: in,
	}

	return value, nil
}

func (t LinodeAutoAllocIPType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t LinodeAutoAllocIPType) ValueType(ctx context.Context) attr.Value {
	return LinodeAutoAllocIPValue{}
}

var _ basetypes.StringValuable = LinodeAutoAllocIPValue{}

// LinodeAutoAllocIPValue represents a string that is either a specific IP address
// or the literal "auto". It implements semantic equality so that a configured
// value of "auto" is considered equal to a concrete IP address stored in state.
type LinodeAutoAllocIPValue struct {
	basetypes.StringValue
}

func (v LinodeAutoAllocIPValue) Equal(o attr.Value) bool {
	other, ok := o.(LinodeAutoAllocIPValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v LinodeAutoAllocIPValue) Type(ctx context.Context) attr.Type {
	return LinodeAutoAllocIPType{}
}

func (v LinodeAutoAllocIPValue) StringSemanticEquals(
	ctx context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(LinodeAutoAllocIPValue)
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

	// A configured value of "auto" is satisfied by any allocated address, so
	// treat it as equal to whatever value is already present to avoid a
	// perpetual diff on subsequent applies.
	if v.ValueString() == "auto" || newValue.ValueString() == "auto" {
		return true, diags
	}

	return v.ValueString() == newValue.ValueString(), diags
}

func LinodeAutoAllocIPValueFrom(value string) LinodeAutoAllocIPValue {
	return LinodeAutoAllocIPValue{
		StringValue: types.StringValue(value),
	}
}
