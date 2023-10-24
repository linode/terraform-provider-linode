//go:build unit

package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestIPAddr_semanticEquals(t *testing.T) {
	v1 := IPAddrStringValue{
		StringValue: types.StringValue("2600:3c02:e000:06b3::"),
	}

	v2 := IPAddrStringValue{
		StringValue: types.StringValue("2600:3c02:e000:6b3::"),
	}

	equal, d := v1.StringSemanticEquals(context.Background(), v2)
	if d.HasError() {
		t.Fatal("Expected no errors; got some")
	}

	if !equal {
		t.Fatal("Expected semantic equality")
	}

	v2.StringValue = types.StringValue("2600:3c02:e000:6c3::")

	equal, d = v1.StringSemanticEquals(context.Background(), v2)
	if d.HasError() {
		t.Fatal("Expected no errors; got some")
	}

	if equal {
		t.Fatal("Expected no semantic equality")
	}
}

func TestIPAddr_validation(t *testing.T) {
	ipType := IPAddrStringType{}

	d := ipType.Validate(context.Background(), tftypes.NewValue(tftypes.String, "192.168.0.1"), path.Empty())
	if d.HasError() {
		t.Fatal("Expected no error; got some")
	}

	d = ipType.Validate(context.Background(), tftypes.NewValue(tftypes.String, "192.1c8.0.1"), path.Empty())
	if !d.HasError() {
		t.Fatal("Expected error; got none")
	}
}
