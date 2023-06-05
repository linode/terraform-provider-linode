package frameworkfilter

import "testing"

func TestResolveStructFieldByJSON(t *testing.T) {
	type TestStruct struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	result, d := resolveStructFieldByJSON(TestStruct{}, "foo")
	if d != nil {
		t.Fatal(d.Detail())
	}

	if result.Name != "Foo" {
		t.Fatalf("Expected Foo; got %s", result.Name)
	}
}

func TestResolveStructValueByJSON(t *testing.T) {
	type TestStruct struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	result, d := resolveStructValueByJSON(
		TestStruct{
			Foo: "cool",
			Bar: "test",
		},
		"foo",
	)
	if d != nil {
		t.Fatal(d.Detail())
	}

	if result.(string) != "cool" {
		t.Fatalf("Expected cool; got %s", result.(string))
	}
}
