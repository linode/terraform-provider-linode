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
	_ basetypes.Int64Typable                    = DomainRecordTTLType{}
	_ basetypes.Int64ValuableWithSemanticEquals = DomainRecordTTLValue{}
)

type DomainRecordTTLType struct {
	basetypes.Int64Type
}

func (t DomainRecordTTLType) Equal(o attr.Type) bool {
	other, ok := o.(DomainRecordTTLType)

	if !ok {
		return false
	}

	return t.Int64Type.Equal(other.Int64Type)
}

func (t DomainRecordTTLType) Int64() string {
	return "DomainRecordTTLType"
}

func (t DomainRecordTTLType) ValueFromInt64(ctx context.Context, in basetypes.Int64Value) (basetypes.Int64Valuable, diag.Diagnostics) {
	value := DomainRecordTTLValue{
		Int64Value: in,
	}

	return value, nil
}

func (t DomainRecordTTLType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.Int64Type.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.Int64Value)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromInt64(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting Int64Value to Int64Valuable: %v", diags)
	}

	return stringValuable, nil
}

func (t DomainRecordTTLType) ValueType(ctx context.Context) attr.Value {
	return DomainRecordTTLValue{}
}

type DomainRecordTTLValue struct {
	basetypes.Int64Value
}

func (v DomainRecordTTLValue) Equal(o attr.Value) bool {
	other, ok := o.(DomainRecordTTLValue)

	if !ok {
		return false
	}

	return v.Int64Value.Equal(other.Int64Value)
}

func (v DomainRecordTTLValue) Type(ctx context.Context) attr.Type {
	return DomainRecordTTLType{}
}

func (v DomainRecordTTLValue) Int64SemanticEquals(ctx context.Context, newValuable basetypes.Int64Valuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(DomainRecordTTLValue)

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

	provisioned := v.Int64Value.ValueInt64()
	declared := newValue.ValueInt64()

	accepted := []int64{
		30, 120, 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, 2419200,
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

	return rounder(declared) == provisioned, diags
}
