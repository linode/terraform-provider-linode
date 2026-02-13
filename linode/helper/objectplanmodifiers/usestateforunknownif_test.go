//go:build unit

package objectplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/objectplanmodifiers"
)

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

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"port": types.Int64Type,
		},
	}

	testCases := map[string]struct {
		request          planmodifier.ObjectRequest
		condition        objectplanmodifiers.UseStateForUnknownIfFunc
		expectedWarnings int
		expectedErrors   int
	}{
		"diagnostics-warning": {
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
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning")
				resp.UseState = false
			},
			expectedWarnings: 1,
			expectedErrors:   0,
		},
		"diagnostics-error": {
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
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 0,
			expectedErrors:   1,
		},
		"diagnostics-multiple": {
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
				resp.Diagnostics.AddWarning("Test Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Test Warning 2", "Second warning")
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 2,
			expectedErrors:   1,
		},
		"diagnostics-none": {
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
			expectedWarnings: 0,
			expectedErrors:   0,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.ObjectResponse{
				PlanValue:   req.PlanValue,
				Diagnostics: diag.Diagnostics{},
			}
			objectplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyObject(context.Background(), req, resp)

			if resp.Diagnostics.WarningsCount() != testCase.expectedWarnings {
				t.Errorf("expected %d warnings, got %d", testCase.expectedWarnings, resp.Diagnostics.WarningsCount())
			}

			if resp.Diagnostics.ErrorsCount() != testCase.expectedErrors {
				t.Errorf("expected %d errors, got %d", testCase.expectedErrors, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}
