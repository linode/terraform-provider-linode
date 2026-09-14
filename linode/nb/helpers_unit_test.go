//go:build unit

package nb

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateIPv4RangeAutoAssignConflict(t *testing.T) {
	t.Parallel()

	t.Run("errors when ipv4_range and ipv4_range_auto_assign are both set", func(t *testing.T) {
		t.Parallel()

		vpcs := types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(1),
				"ipv4_range":             types.StringValue("10.0.0.4/30"),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
		})

		var diags diag.Diagnostics
		validateIPv4RangeAutoAssignConflict(context.Background(), vpcs, "backend_vpcs", &diags)
		assert.True(t, diags.HasError())
		assert.Contains(t, diags.Errors()[0].Detail(), "ipv4_range")
		assert.Contains(t, diags.Errors()[0].Detail(), "ipv4_range_auto_assign")
	})

	t.Run("allows ipv4_range_auto_assign without ipv4_range", func(t *testing.T) {
		t.Parallel()

		vpcs := types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(1),
				"ipv4_range":             types.StringUnknown(),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolValue(true),
			}),
		})

		var diags diag.Diagnostics
		validateIPv4RangeAutoAssignConflict(context.Background(), vpcs, "backend_vpcs", &diags)
		assert.False(t, diags.HasError())
	})

	t.Run("allows ipv4_range without ipv4_range_auto_assign", func(t *testing.T) {
		t.Parallel()

		vpcs := types.ListValueMust(backendVPCObjType, []attr.Value{
			types.ObjectValueMust(backendVPCObjType.AttrTypes, map[string]attr.Value{
				"subnet_id":              types.Int64Value(1),
				"ipv4_range":             types.StringValue("10.0.0.4/30"),
				"ipv6_range":             types.StringNull(),
				"ipv4_range_auto_assign": types.BoolNull(),
			}),
		})

		var diags diag.Diagnostics
		validateIPv4RangeAutoAssignConflict(context.Background(), vpcs, "vpcs", &diags)
		assert.False(t, diags.HasError())
	})

	t.Run("skips null and unknown lists", func(t *testing.T) {
		t.Parallel()

		var diags diag.Diagnostics
		validateIPv4RangeAutoAssignConflict(context.Background(), types.ListNull(backendVPCObjType), "backend_vpcs", &diags)
		validateIPv4RangeAutoAssignConflict(context.Background(), types.ListUnknown(backendVPCObjType), "backend_vpcs", &diags)
		assert.False(t, diags.HasError())
	})
}
