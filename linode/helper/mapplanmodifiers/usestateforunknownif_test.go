//go:build unit

package mapplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/mapplanmodifiers"
)

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.MapRequest
		condition mapplanmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.MapResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
			},
		},
		"custom-condition-non-empty": {
			// custom condition - only use if map is not empty
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("value")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("value")}),
			},
		},
		"custom-condition-empty": {
			// custom condition with empty map - should not use state value
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.MapResponse{
				PlanValue: req.PlanValue,
			}
			mapplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyMap(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	testCases := map[string]struct {
		request          planmodifier.MapRequest
		condition        mapplanmodifiers.UseStateForUnknownIfFunc
		expectedWarnings int
		expectedErrors   int
	}{
		"diagnostics-warning": {
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning")
				resp.UseState = false
			},
			expectedWarnings: 1,
			expectedErrors:   0,
		},
		"diagnostics-error": {
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 0,
			expectedErrors:   1,
		},
		"diagnostics-multiple": {
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Test Warning 2", "Second warning")
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 2,
			expectedErrors:   1,
		},
		"diagnostics-none": {
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.MapRequest, resp *mapplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expectedWarnings: 0,
			expectedErrors:   0,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.MapResponse{
				PlanValue:   req.PlanValue,
				Diagnostics: diag.Diagnostics{},
			}
			mapplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyMap(context.Background(), req, resp)

			if resp.Diagnostics.WarningsCount() != testCase.expectedWarnings {
				t.Errorf("expected %d warnings, got %d", testCase.expectedWarnings, resp.Diagnostics.WarningsCount())
			}

			if resp.Diagnostics.ErrorsCount() != testCase.expectedErrors {
				t.Errorf("expected %d errors, got %d", testCase.expectedErrors, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}
