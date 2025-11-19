//go:build unit

package objectplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/objectplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"port": types.Int64Type,
		},
	}

	testCases := map[string]struct {
		request  planmodifier.ObjectRequest
		expected *planmodifier.ObjectResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.ObjectNull(objectType.AttrTypes),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(objectType.AttrTypes),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("old"),
					"port": types.Int64Value(80),
				}),
				PlanValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("new"),
					"port": types.Int64Value(443),
				}),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("new"),
					"port": types.Int64Value(443),
				}),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectUnknown(objectType.AttrTypes),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(objectType.AttrTypes),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ObjectNull(objectType.AttrTypes),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(objectType.AttrTypes),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.ObjectResponse{
				PlanValue: req.PlanValue,
			}
			objectplanmodifiers.UseStateForUnknownIfNotNull().PlanModifyObject(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf(t *testing.T) {
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"port": types.Int64Type,
		},
	}

	testCases := map[string]struct {
		request   planmodifier.ObjectRequest
		condition objectplanmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.ObjectResponse
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			condition: func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(objectType.AttrTypes),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			condition: func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("state"),
					"port": types.Int64Value(80),
				}),
			},
		},
		"custom-condition-has-name": {
			// custom condition - only use if object has non-empty name attribute
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("valid-name"),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			condition: func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifiers.UseStateForUnknownIfFuncResponse) {
				if req.StateValue.IsNull() {
					resp.UseState = false
					return
				}
				attrs := req.StateValue.Attributes()
				nameAttr, exists := attrs["name"]
				if !exists {
					resp.UseState = false
					return
				}
				nameStr, ok := nameAttr.(types.String)
				if !ok {
					resp.UseState = false
					return
				}
				resp.UseState = !nameStr.IsNull() && nameStr.ValueString() != ""
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue("valid-name"),
					"port": types.Int64Value(80),
				}),
			},
		},
		"custom-condition-empty-name": {
			// custom condition with empty name - should not use state value
			request: planmodifier.ObjectRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue: types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
					"name": types.StringValue(""),
					"port": types.Int64Value(80),
				}),
				PlanValue:   types.ObjectUnknown(objectType.AttrTypes),
				ConfigValue: types.ObjectNull(objectType.AttrTypes),
			},
			condition: func(ctx context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifiers.UseStateForUnknownIfFuncResponse) {
				if req.StateValue.IsNull() {
					resp.UseState = false
					return
				}
				attrs := req.StateValue.Attributes()
				nameAttr, exists := attrs["name"]
				if !exists {
					resp.UseState = false
					return
				}
				nameStr, ok := nameAttr.(types.String)
				if !ok {
					resp.UseState = false
					return
				}
				resp.UseState = !nameStr.IsNull() && nameStr.ValueString() != ""
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(objectType.AttrTypes),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.ObjectResponse{
				PlanValue: req.PlanValue,
			}
			objectplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyObject(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}
