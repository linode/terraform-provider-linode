package boolplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknownIfNotNull is a convenience wrapper to only use the state value
// in place of unknown values in plans if its value is not null.
func UseStateForUnknownIfNotNull() planmodifier.Bool {
	return UseStateForUnknownIf(
		func(ctx context.Context, request planmodifier.BoolRequest, resp *UseStateForUnknownIfFuncResponse) {
			resp.UseState = !request.StateValue.IsNull()
		},
	)
}

type UseStateForUnknownIfFunc func(context.Context, planmodifier.BoolRequest, *UseStateForUnknownIfFuncResponse)

type UseStateForUnknownIfFuncResponse struct {
	// Diagnostics report errors or warnings related to this logic. An empty
	// or unset slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// UseState should be enabled if conditions are met
	UseState bool
}

// UseStateForUnknownIf returns a plan modifier that will only use the state value
// in place of an unknown value if the given condition is met.
func UseStateForUnknownIf(condition UseStateForUnknownIfFunc) planmodifier.Bool {
	return useStateForUnknownIfModifier{
		conditionFunc: condition,
	}
}

type useStateForUnknownIfModifier struct {
	conditionFunc UseStateForUnknownIfFunc
}

func (m useStateForUnknownIfModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change " +
		"if the given condition is met."
}

func (m useStateForUnknownIfModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change " +
		"if the given condition is met."
}

func (m useStateForUnknownIfModifier) PlanModifyBool(
	ctx context.Context,
	req planmodifier.BoolRequest,
	resp *planmodifier.BoolResponse,
) {
	// Do nothing if there is no state (resource is being created).
	if req.State.Raw.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	var ifFuncResp UseStateForUnknownIfFuncResponse
	m.conditionFunc(ctx, req, &ifFuncResp)

	resp.Diagnostics.Append(ifFuncResp.Diagnostics...)

	if !ifFuncResp.UseState {
		return
	}

	resp.PlanValue = req.StateValue
}
