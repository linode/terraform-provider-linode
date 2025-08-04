package customtypes

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringTypable                    = LinodeScopesStringType{}
	_ basetypes.StringValuableWithSemanticEquals = LinodeAutoAllocRangeValue{}
)

type LinodeAutoAllocRangeType struct {
	basetypes.StringType
}

func (t LinodeAutoAllocRangeType) Equal(o attr.Type) bool {
	other, ok := o.(LinodeAutoAllocRangeType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t LinodeAutoAllocRangeType) String() string {
	return "LinodeAutoAllocRangeType"
}

func (t LinodeAutoAllocRangeType) ValueFromString(
	ctx context.Context,
	in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	// LinodeAutoAllocRangeValue defined in the value type section
	value := LinodeAutoAllocRangeValue{
		StringValue: in,
	}

	return value, nil
}

func (t LinodeAutoAllocRangeType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t LinodeAutoAllocRangeType) ValueType(ctx context.Context) attr.Value {
	return LinodeAutoAllocRangeValue{}
}

var _ basetypes.StringValuable = LinodeAutoAllocRangeValue{}

type LinodeAutoAllocRangeValue struct {
	basetypes.StringValue
}

func (v LinodeAutoAllocRangeValue) Equal(o attr.Value) bool {
	other, ok := o.(LinodeAutoAllocRangeValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v LinodeAutoAllocRangeValue) Type(ctx context.Context) attr.Type {
	return LinodeAutoAllocRangeType{}
}

// LinodeAutoAllocRangeValue defined in the value type section
// Ensure the implementation satisfies the expected interfaces

func (v LinodeAutoAllocRangeValue) StringSemanticEquals(
	ctx context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(LinodeAutoAllocRangeValue)
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

	addr, prefix, err := parseRangeOptionalAddress(v.ValueString())
	if err != nil {
		diags.AddError(
			"Failed to Parse Range",
			err.Error(),
		)
		return false, diags
	}

	newAddr, newPrefix, err := parseRangeOptionalAddress(newValue.ValueString())
	if err != nil {
		diags.AddError(
			"Failed to Parse Range",
			err.Error(),
		)
		return false, diags
	}

	// One of the addresses only has a prefix specified,
	// so we should only diff on the prefix.
	if (addr == nil) || (newAddr == nil) {
		return prefix == newPrefix, diags
	}

	return prefix == newPrefix && addr.Compare(*newAddr) == 0, diags
}

func parseRangeOptionalAddress(cidr string) (*netip.Addr, int, error) {
	cidrSplit := strings.Split(cidr, "/")
	if len(cidrSplit) != 2 {
		return nil, 0, fmt.Errorf("malformed CIDR: %s", cidr)
	}

	ipStr, prefixStr := cidrSplit[0], cidrSplit[1]

	if prefixStr == "" {
		return nil, 0, fmt.Errorf("expected non-empty prefix: %s", cidr)
	}

	prefix, err := strconv.Atoi(prefixStr)
	if err != nil {
		return nil, 0, fmt.Errorf("malformed prefix: %w", err)
	}

	// No address is specified
	if ipStr == "" {
		return nil, prefix, nil
	}

	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse CIDR: %w", err)
	}

	return &ip, prefix, nil
}
