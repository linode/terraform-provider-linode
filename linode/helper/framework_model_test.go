package helper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCompareModels(t *testing.T) {
	type TestModel struct {
		Field1 types.Int64  `tfsdk:"field1" linode_mutable:"true"`
		Field2 types.String `tfsdk:"field2" linode_mutable:"true"`
		Field3 types.Bool   `tfsdk:"field3"`
	}

	inst1 := TestModel{
		Field1: types.Int64Value(12345),
		Field2: types.StringValue("bar"),
		Field3: types.BoolValue(true),
	}

	inst2 := TestModel{
		Field1: types.Int64Value(54321),
		Field2: types.StringValue("foo"),
		Field3: types.BoolValue(false),
	}

	result, err := ShouldModelUpdate(&inst1, &inst2)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatal("Expected result to be true; got false")
	}

	result, err = ShouldModelUpdate(inst1, inst1)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("Expected result to be false; got true")
	}

	// Let's make sure field 3 isn't being used for comparisons
	inst2.Field1 = inst1.Field1
	inst2.Field2 = inst1.Field2

	result, err = ShouldModelUpdate(inst1, inst2)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatal("Expected result to be false; got true")
	}
}
