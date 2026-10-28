//go:build unit

package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLinodeAutoAllocIP_semanticEquals(t *testing.T) {
	cases := []struct {
		name  string
		a     string
		b     string
		equal bool
	}{
		{"auto vs resolved IP", "auto", "10.0.1.5", true},
		{"resolved IP vs auto", "10.0.1.5", "auto", true},
		{"auto vs auto", "auto", "auto", true},
		{"same IP", "10.0.1.5", "10.0.1.5", true},
		{"different IP", "10.0.1.5", "10.0.1.6", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v1 := LinodeAutoAllocIPValue{StringValue: types.StringValue(c.a)}
			v2 := LinodeAutoAllocIPValue{StringValue: types.StringValue(c.b)}

			equal, d := v1.StringSemanticEquals(context.Background(), v2)
			if d.HasError() {
				t.Fatalf("Expected no errors; got %v", d)
			}

			if equal != c.equal {
				t.Fatalf("Expected semantic equality %v; got %v", c.equal, equal)
			}
		})
	}
}
