//go:build unit

package natgateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanAddressOps(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		toAdd    []string
		toRemove []string
		headroom int
		expected []addressOp
	}{
		{
			name:     "no changes",
			expected: []addressOp{},
		},
		{
			name:     "addition with headroom",
			toAdd:    []string{"192.0.2.1"},
			headroom: 2,
			expected: []addressOp{{address: "192.0.2.1", add: true}},
		},
		{
			name:     "removal only",
			toRemove: []string{"192.0.2.1"},
			headroom: 0,
			expected: []addressOp{{address: "192.0.2.1"}},
		},
		{
			name:     "replacement with headroom adds before removing",
			toAdd:    []string{"192.0.2.2"},
			toRemove: []string{"192.0.2.1"},
			headroom: 1,
			expected: []addressOp{
				{address: "192.0.2.2", add: true},
				{address: "192.0.2.1"},
			},
		},
		{
			name:     "replacement at the limit removes before adding",
			toAdd:    []string{"192.0.2.2"},
			toRemove: []string{"192.0.2.1"},
			headroom: 0,
			expected: []addressOp{
				{address: "192.0.2.1"},
				{address: "192.0.2.2", add: true},
			},
		},
		{
			name:     "multiple replacements at the limit interleave",
			toAdd:    []string{"192.0.2.3", "192.0.2.4"},
			toRemove: []string{"192.0.2.1", "192.0.2.2"},
			headroom: 0,
			expected: []addressOp{
				{address: "192.0.2.1"},
				{address: "192.0.2.3", add: true},
				{address: "192.0.2.2"},
				{address: "192.0.2.4", add: true},
			},
		},
		{
			name:     "partial headroom adds first then interleaves",
			toAdd:    []string{"192.0.2.3", "192.0.2.4"},
			toRemove: []string{"192.0.2.1", "192.0.2.2"},
			headroom: 1,
			expected: []addressOp{
				{address: "192.0.2.3", add: true},
				{address: "192.0.2.1"},
				{address: "192.0.2.4", add: true},
				{address: "192.0.2.2"},
			},
		},
		{
			name:     "addition without headroom is still attempted",
			toAdd:    []string{"192.0.2.1"},
			headroom: 0,
			expected: []addressOp{{address: "192.0.2.1", add: true}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, planAddressOps(tc.toAdd, tc.toRemove, tc.headroom))
		})
	}
}

// TestPlanAddressOps_neverExceedsLimit verifies that applying the planned
// operations in order never pushes the assigned address count above the limit.
func TestPlanAddressOps_neverExceedsLimit(t *testing.T) {
	t.Parallel()

	const limit = 2

	assigned := 2
	ops := planAddressOps(
		[]string{"192.0.2.3", "192.0.2.4"},
		[]string{"192.0.2.1", "192.0.2.2"},
		limit-assigned,
	)

	for _, op := range ops {
		if op.add {
			assigned++
		} else {
			assigned--
		}

		assert.LessOrEqual(t, assigned, limit, "assigned addresses exceeded the limit")
		assert.GreaterOrEqual(t, assigned, 1, "gateway was left without an address")
	}

	assert.Equal(t, 2, assigned)
}

func TestAddressesMissingFrom(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		[]string{"192.0.2.1"},
		addressesMissingFrom([]string{"192.0.2.1", "192.0.2.2"}, []string{"192.0.2.2"}),
	)
	assert.Empty(t, addressesMissingFrom(nil, []string{"192.0.2.1"}))
	assert.Equal(t, []string{"192.0.2.1"}, addressesMissingFrom([]string{"192.0.2.1"}, nil))
}
