//go:build unit

package stringplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/stringplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.StringRequest
		expected *planmodifier.StringResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("old-value"),
				PlanValue:   types.StringValue("new-value"),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("new-value"),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringUnknown(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("state-value"),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.StringResponse{
				PlanValue: req.PlanValue,
			}
			stringplanmodifiers.UseStateForUnknownIfNotNull().PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.StringRequest
		condition stringplanmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.StringResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("state-value"),
			},
		},
		"custom-condition": {
			// custom condition - only use if value is not empty
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("non-empty"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueString() != ""
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("non-empty"),
			},
		},
		"custom-condition-empty": {
			// custom condition with empty string - should not use state value
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue(""),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueString() != ""
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.StringResponse{
				PlanValue: req.PlanValue,
			}
			stringplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	testCases := map[string]struct {
		request             planmodifier.StringRequest
		condition           stringplanmodifiers.UseStateForUnknownIfFunc
		expectedPlan        types.String
		expectedDiagCount   int
		expectedDiagSummary string
	}{
		"condition-adds-warning": {
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning from condition")
				resp.UseState = true
			},
			expectedPlan:        types.StringValue("state-value"),
			expectedDiagCount:   1,
			expectedDiagSummary: "Test Warning",
		},
		"condition-adds-error": {
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddError("Test Error", "This is a test error from condition")
				resp.UseState = false
			},
			expectedPlan:        types.StringUnknown(),
			expectedDiagCount:   1,
			expectedDiagSummary: "Test Error",
		},
		"condition-adds-multiple-diagnostics": {
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Warning 2", "Second warning")
				resp.UseState = true
			},
			expectedPlan:      types.StringValue("state-value"),
			expectedDiagCount: 2,
		},
		"no-diagnostics": {
			request: planmodifier.StringRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.StringValue("state-value"),
				PlanValue:   types.StringUnknown(),
				ConfigValue: types.StringNull(),
			},
			condition: func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expectedPlan:      types.StringValue("state-value"),
			expectedDiagCount: 0,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.StringResponse{
				PlanValue: req.PlanValue,
			}
			stringplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expectedPlan) {
				t.Errorf("expected plan value %s, got %s", testCase.expectedPlan, resp.PlanValue)
			}

			if len(resp.Diagnostics) != testCase.expectedDiagCount {
				t.Errorf("expected %d diagnostics, got %d", testCase.expectedDiagCount, len(resp.Diagnostics))
			}

			if testCase.expectedDiagSummary != "" && len(resp.Diagnostics) > 0 {
				if resp.Diagnostics[0].Summary() != testCase.expectedDiagSummary {
					t.Errorf("expected diagnostic summary %q, got %q", testCase.expectedDiagSummary, resp.Diagnostics[0].Summary())
				}
			}
		})
	}
}
