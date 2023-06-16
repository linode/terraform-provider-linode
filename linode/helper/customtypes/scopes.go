package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringTypable = LinodeScopesStringType{}
var _ xattr.TypeWithValidate = LinodeScopesStringType{}
var _ basetypes.StringValuableWithSemanticEquals = LinodeScopesStringValue{}

var AvailableScopes []string = []string{
	"account:read_only",
	"account:read_write",
	"databases:read_only",
	"databases:read_write",
	"domains:read_only",
	"domains:read_write",
	"events:read_only",
	"events:read_write",
	"firewall:read_only",
	"firewall:read_write",
	"images:read_only",
	"images:read_write",
	"ips:read_only",
	"ips:read_write",
	"linodes:read_only",
	"linodes:read_write",
	"lke:read_only",
	"lke:read_write",
	"longview:read_only",
	"longview:read_write",
	"nodebalancers:read_only",
	"nodebalancers:read_write",
	"object_storage:read_only",
	"object_storage:read_write",
	"stackscripts:read_only",
	"stackscripts:read_write",
	"volumes:read_only",
	"volumes:read_write",
}

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

func (t LinodeScopesStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
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

func (v LinodeScopesStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
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

// LinodeScopesStringType defined in the schema type section
func (t LinodeScopesStringType) Validate(ctx context.Context, value tftypes.Value, valuePath path.Path) diag.Diagnostics {
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

	if valueString == "*" {
		return diags
	}

	scopes := strings.Split(valueString, " ")

	for _, declaredScope := range scopes {
		validDeclaredScope := false
		for _, availableScope := range AvailableScopes {
			if declaredScope == availableScope {
				validDeclaredScope = true
			}
		}
		if !validDeclaredScope {
			diags.AddAttributeError(
				valuePath,
				"Invalid Scope",
				fmt.Sprintf("'%s' is not a valid scope.", declaredScope),
			)
			return diags
		}
	}

	return diags
}
