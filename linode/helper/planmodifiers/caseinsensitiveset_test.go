//go:build unit

package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCaseInsensitiveSet(t *testing.T) {
	testCases := []struct {
		Old, New, Expected []string
	}{
		{
			Old:      []string{"foo", "bar"},
			New:      []string{"foo", "bar"},
			Expected: []string{"foo", "bar"},
		},
		{
			Old:      []string{"foo", "Bar"},
			New:      []string{"foo", "bar"},
			Expected: []string{"foo", "Bar"},
		},
		{
			Old:      []string{"foo", "bar"},
			New:      []string{"fOO", "bar"},
			Expected: []string{"foo", "bar"},
		},
		{
			Old:      []string{"foo", "bar"},
			New:      []string{"fOO", "bar", "wow"},
			Expected: []string{"foo", "bar", "wow"},
		},
		{
			Old:      []string{"foo", "bar", "wOw"},
			New:      []string{"fOO", "bar"},
			Expected: []string{"foo", "bar"},
		},
	}

	var diags diag.Diagnostics

	testPlanModifier := CaseInsensitiveSet()

	for i, testCase := range testCases {
		stateValue, d := types.SetValueFrom(
			context.Background(),
			types.StringType,
			testCase.Old,
		)
		diags.Append(d...)

		planValue, d := types.SetValueFrom(
			context.Background(),
			types.StringType,
			testCase.New,
		)
		diags.Append(d...)

		expectedValue, d := types.SetValueFrom(
			context.Background(),
			types.StringType,
			testCase.Expected,
		)
		diags.Append(d...)

		if d.HasError() {
			t.Fatalf("%d: got error building sets: %v", i, d.Errors())
		}

		req := planmodifier.SetRequest{
			StateValue: stateValue,
			PlanValue:  planValue,
		}

		var resp planmodifier.SetResponse

		testPlanModifier.PlanModifySet(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("%d: got error modifying plan: %v", i, resp.Diagnostics.Errors())
		}

		if !resp.PlanValue.Equal(expectedValue) {
			t.Fatalf("%d: output plan value does not equal expected plan value", i)
		}
	}
}
