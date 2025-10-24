//go:build unit

package setplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/setplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.SetRequest
		expected *planmodifier.SetResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("old")}),
				PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("new")}),
				ConfigValue: types.SetNull(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("new")}),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetUnknown(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.SetResponse{
				PlanValue: req.PlanValue,
			}
			setplanmodifiers.UseStateForUnknownIfNotNull().PlanModifySet(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.SetRequest
		condition setplanmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.SetResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
			},
		},
		"custom-condition-non-empty": {
			// custom condition - only use if set is not empty
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{types.StringValue("item")}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("item")}),
			},
		},
		"custom-condition-empty": {
			// custom condition with empty set - should not use state value
			request: planmodifier.SetRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.SetValueMust(types.StringType, []attr.Value{}),
				PlanValue:   types.SetUnknown(types.StringType),
				ConfigValue: types.SetNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.SetResponse{
				PlanValue: types.SetUnknown(types.StringType),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.SetResponse{
				PlanValue: req.PlanValue,
			}
			setplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifySet(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}
