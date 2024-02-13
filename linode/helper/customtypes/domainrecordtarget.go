package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = DomainRecordTargetType{}
	_ basetypes.StringValuableWithSemanticEquals = DomainRecordTargetValue{}
)

type DomainRecordTargetType struct {
	basetypes.StringType
}

func (t DomainRecordTargetType) Equal(o attr.Type) bool {
	other, ok := o.(DomainRecordTargetType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t DomainRecordTargetType) String() string {
	return "DomainRecordTargetType"
}

func (t DomainRecordTargetType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := DomainRecordTargetValue{
		StringValue: in,
	}

	return value, nil
}

func (t DomainRecordTargetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t DomainRecordTargetType) ValueType(ctx context.Context) attr.Value {
	return DomainRecordTargetValue{}
}

type DomainRecordTargetValue struct {
	basetypes.StringValue
}

func (v DomainRecordTargetValue) Equal(o attr.Value) bool {
	other, ok := o.(DomainRecordTargetValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v DomainRecordTargetValue) Type(ctx context.Context) attr.Type {
	return DomainRecordTargetType{}
}

func (v DomainRecordTargetValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(DomainRecordTargetValue)

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

	provisioned := v.StringValue.ValueString()
	declared := newValue.ValueString()

	return len(strings.Split(declared, ".")) == 1 &&
		strings.Contains(provisioned, declared), diags
}
