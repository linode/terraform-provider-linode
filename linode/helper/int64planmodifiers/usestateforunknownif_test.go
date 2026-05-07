//go:build unit

package int64planmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/int64planmodifiers"
)

func TestUseStateForUnknownIf(t *testing.T) {
	testCases := map[string]struct {
		request   planmodifier.Int64Request
		condition int64planmodifiers.UseStateForUnknownIfFunc
		expected  *planmodifier.Int64Response
	}{
		"condition-false": {
			// condition returns false, should not use state value
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = false
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Unknown(),
			},
		},
		"condition-true": {
			// condition returns true, should use state value
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = true
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(42),
			},
		},
		"custom-condition-greater-than-zero": {
			// custom condition - only use if value is greater than zero
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(100),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueInt64() > 0
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Value(100),
			},
		},
		"custom-condition-zero": {
			// custom condition with zero - should not use state value
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(0),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.UseState = !req.StateValue.IsNull() && req.StateValue.ValueInt64() > 0
			},
			expected: &planmodifier.Int64Response{
				PlanValue: types.Int64Unknown(),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.Int64Response{
				PlanValue: req.PlanValue,
			}
			int64planmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyInt64(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}

func TestUseStateForUnknownIf_Diagnostics(t *testing.T) {
	testCases := map[string]struct {
		request          planmodifier.Int64Request
		condition        int64planmodifiers.UseStateForUnknownIfFunc
		expectedWarnings int
		expectedErrors   int
	}{
		"diagnostics-warning": {
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning", "This is a test warning")
				resp.UseState = false
			},
			expectedWarnings: 1,
			expectedErrors:   0,
		},
		"diagnostics-error": {
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 0,
			expectedErrors:   1,
		},
		"diagnostics-multiple": {
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
				resp.Diagnostics.AddWarning("Test Warning 1", "First warning")
				resp.Diagnostics.AddWarning("Test Warning 2", "Second warning")
				resp.Diagnostics.AddError("Test Error", "This is a test error")
				resp.UseState = false
			},
			expectedWarnings: 2,
			expectedErrors:   1,
		},
		"diagnostics-none": {
			request: planmodifier.Int64Request{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.Int64Value(42),
				PlanValue:   types.Int64Unknown(),
				ConfigValue: types.Int64Null(),
			},
			condition: func(ctx context.Context, req planmodifier.Int64Request, resp *int64planmodifiers.UseStateForUnknownIfFuncResponse) {
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
			resp := &planmodifier.Int64Response{
				PlanValue:   req.PlanValue,
				Diagnostics: diag.Diagnostics{},
			}
			int64planmodifiers.UseStateForUnknownIf(testCase.condition).PlanModifyInt64(context.Background(), req, resp)

			if resp.Diagnostics.WarningsCount() != testCase.expectedWarnings {
				t.Errorf("expected %d warnings, got %d", testCase.expectedWarnings, resp.Diagnostics.WarningsCount())
			}

			if resp.Diagnostics.ErrorsCount() != testCase.expectedErrors {
				t.Errorf("expected %d errors, got %d", testCase.expectedErrors, resp.Diagnostics.ErrorsCount())
			}
		})
	}
}
