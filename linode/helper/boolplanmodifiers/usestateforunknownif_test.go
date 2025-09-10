//go:build unit

package boolplanmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/boolplanmodifiers"
)

func TestUseStateForUnknownIfNotNull(t *testing.T) {
	testCases := map[string]struct {
		request  planmodifier.BoolRequest
		expected *planmodifier.BoolResponse
	}{
		"null-state": {
			// resource creation - state is null
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, nil),
				},
				StateValue:  types.BoolNull(),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolUnknown(),
			},
		},
		"known-plan": {
			// the plan is already known, don't change it
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolValue(false),
				ConfigValue: types.BoolNull(),
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolValue(false),
			},
		},
		"unknown-config": {
			// the config is unknown, don't interfere
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolUnknown(),
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolUnknown(),
			},
		},
		"null-state-value": {
			// the state value is null, don't use it
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolNull(),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolUnknown(),
			},
		},
		"use-state-value": {
			// should use the state value
			request: planmodifier.BoolRequest{
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
				},
				StateValue:  types.BoolValue(true),
				PlanValue:   types.BoolUnknown(),
				ConfigValue: types.BoolNull(),
			},
			expected: &planmodifier.BoolResponse{
				PlanValue: types.BoolValue(true),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := testCase.request
			resp := &planmodifier.BoolResponse{
				PlanValue: req.PlanValue,
			}
			boolplanmodifiers.UseStateForUnknownIfNotNull().PlanModifyBool(context.Background(), req, resp)

			if !resp.PlanValue.Equal(testCase.expected.PlanValue) {
				t.Errorf("expected %s, got %s", testCase.expected.PlanValue, resp.PlanValue)
			}
		})
	}
}
