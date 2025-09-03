//go:build unit

package mapplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/mapplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.MapRequest
		expected *planmodifier.MapResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.MapNull(types.StringType),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("old")}),
				PlanValue:   types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("new")}),
				ConfigValue: types.MapNull(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("new")}),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapUnknown(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapNull(types.StringType),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapUnknown(types.StringType),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.MapRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
				PlanValue:   types.MapUnknown(types.StringType),
				ConfigValue: types.MapNull(types.StringType),
			},
			expected: &planmodifier.MapResponse{
				PlanValue: types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue("state")}),
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
			mapplanmodifiers.UseStateForUnknownIfNotNull().PlanModifyMap(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.MapRequest
		condition func(context.Context, planmodifier.MapRequest) bool
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
			condition: func(ctx context.Context, req planmodifier.MapRequest) bool {
				return false
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
			condition: func(ctx context.Context, req planmodifier.MapRequest) bool {
				return true
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
			condition: func(ctx context.Context, req planmodifier.MapRequest) bool {
				elements := req.StateValue.Elements()
				return !req.StateValue.IsNull() && len(elements) > 0
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
			condition: func(ctx context.Context, req planmodifier.MapRequest) bool {
				elements := req.StateValue.Elements()
				return !req.StateValue.IsNull() && len(elements) > 0
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
