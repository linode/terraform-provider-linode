//go:build unit

package listplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/listplanmodifiers"
)

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.ListRequest
		condition listplanmodifiers.UseStateForUnknownIfFunc
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
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
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
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
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
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
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
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				elements := req.StateValue.Elements()
				resp.UseState = !req.StateValue.IsNull() && len(elements) > 0
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

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	testCases := map[string]struct {
		request          planmodifier.ListRequest
		condition        listplanmodifiers.UseStateForUnknownIfFunc
		expectedWarnings int
		expectedErrors   int
	}{
		"diagnostics-warning": {
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning")
				resp.UseState = false
			},
			expectedWarnings: 1,
			expectedErrors:   0,
		},
		"diagnostics-error": {
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 0,
			expectedErrors:   1,
		},
		"diagnostics-multiple": {
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Test Warning 2", "Second warning")
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 2,
			expectedErrors:   1,
		},
		"diagnostics-none": {
			request: planmodifier.ListRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("state")}),
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			},
			condition: func(ctx context.Context, req planmodifier.ListRequest, resp *listplanmodifiers.UseStateForUnknownIfFuncResponse) {
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
			resp := &planmodifier.ListResponse{
				PlanValue:   req.PlanValue,
				Diagnostics: diag.Diagnostics{},
			}
			listplanmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyList(context.Background(), req, resp)

			if resp.Diagnostics.WarningsCount() != testCase.expectedWarnings {
				t.Errorf("expected %d warnings, got %d", testCase.expectedWarnings, resp.Diagnostics.WarningsCount())
			}

			if resp.Diagnostics.ErrorsCount() != testCase.expectedErrors {
				t.Errorf("expected %d errors, got %d", testCase.expectedErrors, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}
