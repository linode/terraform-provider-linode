//go:build unit

package float64planmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/float64planmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.Float64Request
		expected *planmodifier.Float64Response
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.Float64Null(),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Value(45.6),
				ConfigValue: types.Float64Null(),
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Value(45.6),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Unknown(),
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Null(),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Value(12.3),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.Float64Response{
				PlanValue: req.PlanValue,
			}
			float64planmodifiers.UseStateForUnknownIfNotNull().PlanModifyFloat64(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.Float64Request
		condition float64planmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.Float64Response
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Value(12.3),
			},
		},
		"custom-condition-greater-than-zero": {
			// custom condition - only use if value is greater than 0
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(12.3),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueFloat64() > 0
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Value(12.3),
			},
		},
		"custom-condition-zero": {
			// custom condition with zero value - should not use state value
			request: planmodifier.Float64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Float64Value(0.0),
				PlanValue:   types.Float64Unknown(),
				ConfigValue: types.Float64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueFloat64() > 0
			},
			expected: &planmodifier.Float64Response{
				PlanValue: types.Float64Unknown(),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.Float64Response{
				PlanValue: req.PlanValue,
			}
			float64planmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyFloat64(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}
