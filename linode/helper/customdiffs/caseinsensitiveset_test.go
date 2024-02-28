//go:build unit

package customdiffs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCaseInsensitiveSet(t *testing.T) {
	testCases := []struct {
		Old, New, Expected []any
	}{
		{
			Old:      []any{"foo", "bar"},
			New:      []any{"foo", "bar"},
			Expected: []any{"foo", "bar"},
		},
		{
			Old:      []any{"foo", "Bar"},
			New:      []any{"foo", "bar"},
			Expected: []any{"foo", "Bar"},
		},
		{
			Old:      []any{"foo", "bar"},
			New:      []any{"fOO", "bar"},
			Expected: []any{"foo", "bar"},
		},
		{
			Old:      []any{"foo", "bar"},
			New:      []any{"fOO", "bar", "wow"},
			Expected: []any{"foo", "bar", "wow"},
		},
		{
			Old:      []any{"foo", "bar", "wOw"},
			New:      []any{"fOO", "bar"},
			Expected: []any{"foo", "bar"},
		},
	}

	for i, testCase := range testCases {
		expectedResult := schema.NewSet(schema.HashString, testCase.Expected)

		result := computeCaseInsensitivePlannedSet(
			schema.NewSet(schema.HashString, testCase.Old),
			schema.NewSet(schema.HashString, testCase.New),
		)

		if !result.Equal(expectedResult) {
			t.Fatalf("expected result mismatch on case %d", i)
		}
	}
}
