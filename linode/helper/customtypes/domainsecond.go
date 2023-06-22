package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.Int64Typable                    = LinodeDomainSecondsType{}
	_ basetypes.Int64ValuableWithSemanticEquals = LinodeDomainSecondsValue{}
)

type LinodeDomainSecondsType struct {
	basetypes.Int64Type
}

func (t LinodeDomainSecondsType) Equal(o attr.Type) bool {
	other, ok := o.(LinodeDomainSecondsType)
	if !ok {
		return false
	}

	return t.Int64Type.Equal(other.Int64Type)
}

func (t LinodeDomainSecondsType) String() string {
	return "LinodeDomainSecondsType"
}

func (t LinodeDomainSecondsType) ValueFromInt64(
	ctx context.Context,
	in basetypes.Int64Value,
) (basetypes.Int64Valuable, diag.Diagnostics) {
	value := LinodeDomainSecondsValue{
		Int64Value: in,
	}

	return value, nil
}

func (t LinodeDomainSecondsType) ValueFromTerraform(
	ctx context.Context,
	in tftypes.Value,
) (attr.Value, error) {
	attrValue, err := t.Int64Type.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	int64Value, ok := attrValue.(basetypes.Int64Value)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	int64Valuable, diags := t.ValueFromInt64(ctx, int64Value)

	if diags.HasError() {
		return nil, fmt.Errorf(
			"unexpected error converting Int64Value to Int64Valuable: %v",
			diags,
		)
	}

	return int64Valuable, nil
}

func (t LinodeDomainSecondsType) ValueType(ctx context.Context) attr.Value {
	return LinodeDomainSecondsValue{}
}

var _ basetypes.Int64Valuable = LinodeDomainSecondsValue{}

type LinodeDomainSecondsValue struct {
	basetypes.Int64Value
}

func (v LinodeDomainSecondsValue) Equal(o attr.Value) bool {
	other, ok := o.(LinodeDomainSecondsValue)
	if !ok {
		return false
	}

	return v.Int64Value.Equal(other.Int64Value)
}

func (v LinodeDomainSecondsValue) Type(ctx context.Context) attr.Type {
	return LinodeDomainSecondsType{}
}

func (v LinodeDomainSecondsValue) Int64SemanticEquals(
	ctx context.Context,
	newValuable basetypes.Int64Valuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(LinodeDomainSecondsValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value types was received while performing semantic equality checks. "+
				"Please report this to the providers developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	return compareDomainSeconds(v.Int64Value.ValueInt64(), newValue.ValueInt64()), diags
}

func compareDomainSeconds(oldValue, newValue int64) bool {
	accepted := []int64{
		30, 120, 300, 3600, 7200, 14400, 28800, 57600,
		86400, 172800, 345600, 604800, 1209600, 2419200,
	}

	rounder := func(n int64) int64 {
		if n == 0 {
			return 0
		}

		for _, value := range accepted {
			if n <= value {
				return value
			}
		}
		return accepted[len(accepted)-1]
	}

	return rounder(newValue) == oldValue
}
