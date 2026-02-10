package boolplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/boolplanmodifiers"
)

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.BoolRequest
		condition boolplanmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.BoolResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolUnknown(),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolValue(true),
			},
		},
		"custom-condition-true-value": {
			// custom condition - only use if value is true
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueBool()
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolValue(true),
			},
		},
		"custom-condition-false-value": {
			// custom condition with false - should not use state value
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(false),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueBool()
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolUnknown(),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.BoolResponse{
				PlanValue: req.PlanValue,
			}
			boolplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyBool(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	testCases := map[string]struct {
		request          planmodifier.BoolRequest
		condition        boolplanmodifiers.UseStateForUnknownIfFunc
		expectedWarnings int
		expectedErrors   int
	}{
		"diagnostics-warning": {
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning")
				resp.UseState = false
			},
			expectedWarnings: 1,
			expectedErrors:   0,
		},
		"diagnostics-error": {
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 0,
			expectedErrors:   1,
		},
		"diagnostics-multiple": {
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Test Warning 2", "Second warning")
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 2,
			expectedErrors:   1,
		},
		"diagnostics-none": {
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			condition: func(ctx context.Context, req planmodifier.BoolRequest, resp *boolplanmodifiers.UseStateForUnknownIfFuncResponse) {
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
			resp := &planmodifier.BoolResponse{
				PlanValue:   req.PlanValue,
				Diagnostics: diag.Diagnostics{},
			}
			boolplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyBool(context.Background(), req, resp)

			if resp.Diagnostics.WarningsCount() != testCase.expectedWarnings {
				t.Errorf("expected %d warnings, got %d", testCase.expectedWarnings, resp.Diagnostics.WarningsCount())
			}

			if resp.Diagnostics.ErrorsCount() != testCase.expectedErrors {
				t.Errorf("expected %d errors, got %d", testCase.expectedErrors, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}
