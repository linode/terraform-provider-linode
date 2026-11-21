//go:build unit

package helper_test

import (
	"slices"
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func TestCompareStringSets(t *testing.T) {
	cases := []struct {
		name string
		a, b []string
		want bool
	}{
		{name: "both empty", a: nil, b: []string{}, want: true},
		{name: "same order", a: []string{"a", "b"}, b: []string{"a", "b"}, want: true},
		{name: "different order", a: []string{"b", "a"}, b: []string{"a", "b"}, want: true},
		{name: "second unsorted", a: []string{"a", "b", "c"}, b: []string{"c", "a", "b"}, want: true},
		{name: "different items", a: []string{"a", "b"}, b: []string{"a", "c"}, want: false},
		{name: "empty and two items", a: nil, b: []string{"x", "y"}, want: false},
		{name: "one and two items", a: []string{"x"}, b: []string{"x", "y"}, want: false},
		{name: "two and one item", a: []string{"x", "y"}, b: []string{"x"}, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := helper.CompareStringSets(tc.a, tc.b); got != tc.want {
				t.Fatalf("CompareStringSets(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestCompareSetsDoesNotReorderInput(t *testing.T) {
	a := []any{"b", "a"}
	b := []any{"a", "b"}

	if !helper.CompareSets(a, b) {
		t.Fatal("expected the sets to be equal")
	}
	if !slices.Equal(a, []any{"b", "a"}) || !slices.Equal(b, []any{"a", "b"}) {
		t.Fatalf("inputs were reordered: a=%v b=%v", a, b)
	}
}
