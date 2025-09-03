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
		condition func(context.Context, planmodifier.StringRequest) bool
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
			condition: func(ctx context.Context, req planmodifier.StringRequest) bool {
				return false
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
			condition: func(ctx context.Context, req planmodifier.StringRequest) bool {
				return true
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
			condition: func(ctx context.Context, req planmodifier.StringRequest) bool {
				return !req.StateValue.IsNull() && req.StateValue.ValueString() != ""
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
			condition: func(ctx context.Context, req planmodifier.StringRequest) bool {
				return !req.StateValue.IsNull() && req.StateValue.ValueString() != ""
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
