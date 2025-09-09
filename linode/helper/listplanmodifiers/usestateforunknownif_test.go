//go:build unit

package listplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/listplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.ListRequest
		expected *planmodifier.ListResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.ListNull(types.StringType),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListUnknown(types.StringType),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("old")}),
				PlanValue:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("new")}),
				ConfigValue: types.ListNull(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("new")}),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListUnknown(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListUnknown(types.StringType),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListNull(types.StringType),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListUnknown(types.StringType),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.ListResponse{
				PlanValue: req.PlanValue,
			}
			listplanmodifiers.UseStateForUnknownIfNotNull().PlanModifyList(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.ListRequest
		condition func(context.Context, planmodifier.ListRequest) bool
		expected  *planmodifier.ListResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest) bool {
				return false
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListUnknown(types.StringType),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest) bool {
				return true
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
			},
		},
		"custom-condition-non-empty": {
			// custom condition - only use if list is not empty
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("item")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest) bool {
				elements := req.StateValue.Elements()
				return !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("item")}),
			},
		},
		"custom-condition-empty": {
			// custom condition with empty list - should not use state value
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest) bool {
				elements := req.StateValue.Elements()
				return !req.StateValue.IsNull() && len(elements) > 0
			},
			expected: &planmodifier.ListResponse{
				PlanValue: types.ListUnknown(types.StringType),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.ListResponse{
				PlanValue: req.PlanValue,
			}
			listplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyList(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}
