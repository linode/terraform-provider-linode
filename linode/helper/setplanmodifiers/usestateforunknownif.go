package setplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknownUnlessTheseChanged is a convenience wrapper to only use the state value
// in place of unknown values in plans unless another attribute has been changed.
func UseStateForUnknownUnlessTheseChanged(expressions ...path.Expression) planmodifier.Set {
	return UseStateForUnknownIf(
		func(ctx context.Context, request planmodifier.SetRequest, resp *UseStateForUnknownIfFuncResponse) {
			if len(expressions) == 0 {
				resp.UseState = true
				return
			}

			expressions := request.PathExpression.MergeExpressions(expressions...)

			for _, expression := range expressions {
				matchedPaths, newDiags := request.Config.PathMatches(ctx, expression)

				resp.Diagnostics.Append(newDiags...)
				if resp.Diagnostics.HasError() {
					return
				}

				for _, mp := range matchedPaths {
					var state, plan attr.Value

					newDiags = request.Plan.GetAttribute(ctx, mp, &plan)

					resp.Diagnostics.Append(newDiags...)
					if resp.Diagnostics.HasError() {
						return
					}

					if plan.IsUnknown() {
						continue
					}

					newDiags = request.State.GetAttribute(ctx, mp, &state)

					resp.Diagnostics.Append(newDiags...)
					if resp.Diagnostics.HasError() {
						return
					}

					if !state.Equal(plan) {
						// panic(fmt.Sprintf("state and config should not be equal here: state=%#v config=%#v, state_is_null=%v, config_is_null=%v", state, config, state.IsNull(), config.IsNull()))
						resp.UseState = false
						return
					}
				}
			}

			resp.UseState = true
		},
	)
}

// UseStateForUnknownIfNotNull is a convenience wrapper to only use the state value
// in place of unknown values in plans if its value is not null.
func UseStateForUnknownIfNotNull() planmodifier.Set {
	return UseStateForUnknownIf(
		func(ctx context.Context, request planmodifier.SetRequest, resp *UseStateForUnknownIfFuncResponse) {
			resp.UseState = !request.StateValue.IsNull()
		},
	)
}

type UseStateForUnknownIfFunc func(context.Context, planmodifier.SetRequest, *UseStateForUnknownIfFuncResponse)

type UseStateForUnknownIfFuncResponse struct {
	// Diagnostics report errors or warnings related to this logic. An empty
	// or unset slice indicates success, with no warnings or errors generated.
	Diagnostics diag.Diagnostics

	// UseState should be enabled if conditions are met
	UseState bool
}

// UseStateForUnknownIf returns a plan modifier that will only use the state value
// in place of an unknown value if the given condition is met.
func UseStateForUnknownIf(condition UseStateForUnknownIfFunc) planmodifier.Set {
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

func (m useStateForUnknownIfModifier) PlanModifySet(
	ctx context.Context,
	req planmodifier.SetRequest,
	resp *planmodifier.SetResponse,
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

	ifFuncResp := &UseStateForUnknownIfFuncResponse{}

	if m.conditionFunc(ctx, req, ifFuncResp); !ifFuncResp.UseState {
		return
	}

	resp.PlanValue = req.StateValue
}
