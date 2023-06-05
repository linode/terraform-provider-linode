package frameworkfilter

import (
	"log"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestApplyLocalFiltering(t *testing.T) {
	type FilterableStruct struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	testFiltersModel := []FilterModel{
		{
			Name: types.StringValue("foo"),
			Values: []types.String{
				types.StringValue("ba.*"),
				types.StringValue("te.*"),
			},
			MatchBy: types.StringValue("regex"),
		},
		{
			Name: types.StringValue("bar"),
			Values: []types.String{
				types.StringValue("test"),
			},
		},
		{
			Name: types.StringValue("api_foo"),
			Values: []types.String{
				types.StringValue("wow"),
			},
		},
	}

	filterEntries := []FilterableStruct{
		{
			Foo: "test",
			Bar: "test1",
		},
		{
			Foo: "foo",
			Bar: "test",
		},
		{
			Foo: "bar",
			Bar: "test",
		},
	}

	result, d := testFilterConfig.applyLocalFiltering(
		testFiltersModel,
		helper.TypedSliceToAny(filterEntries),
	)
	if d != nil {
		t.Fatal(d.Detail())
	}

	if !reflect.DeepEqual(result[0], filterEntries[2]) {
		t.Fatal(cmp.Diff(result[0], filterEntries[2]))
	}
}

func TestNormalizeValue(t *testing.T) {
	type entry struct {
		ExpectedOutput string
		Input          any
	}
	testCases := []entry{
		{
			"123",
			123,
		},
		{
			"true",
			true,
		},
		{
			"blah",
			"blah",
		},
		{
			"12346",
			12345.678,
		},
	}

	for _, entry := range testCases {
		result, d := normalizeValue(entry.Input)
		if d != nil {
			t.Fatal(d.Detail())
		}

		if result != entry.ExpectedOutput {
			t.Fatalf("%s != %s", result, entry.ExpectedOutput)
		}
	}

}

func TestCheckFilterRegex(t *testing.T) {
	result, d := checkFilterRegex(
		[]types.String{
			types.StringValue("cool.*"),
			types.StringValue("bl.*"),
		},
		"blah",
	)
	if d != nil {
		log.Fatal(d.Detail())
	}
	if !result {
		t.Fatal("Expected true, got false")
	}

	result, d = checkFilterRegex(
		[]types.String{
			types.StringValue("no"), types.StringValue("bad"),
		},
		"blah",
	)
	if d != nil {
		log.Fatal(d.Detail())
	}
	if result {
		t.Fatal("Expected false, got true")
	}
}

func TestCheckFilterSubString(t *testing.T) {
	result, d := checkFilterSubString(
		[]types.String{
			types.StringValue("bl"),
		},
		"blah",
	)
	if d != nil {
		log.Fatal(d.Detail())
	}
	if !result {
		t.Fatal("Expected true, got false")
	}

	result, d = checkFilterSubString(
		[]types.String{
			types.StringValue("no"), types.StringValue("bad"),
		},
		"blah",
	)
	if d != nil {
		log.Fatal(d.Detail())
	}
	if result {
		t.Fatal("Expected false, got true")
	}
}

func TestCheckFilterExact(t *testing.T) {
	result := checkFilterExact(
		[]types.String{
			types.StringValue("blah"),
		},
		"blah",
	)
	if !result {
		t.Fatal("Expected true, got false")
	}

	result = checkFilterExact(
		[]types.String{
			types.StringValue("no"), types.StringValue("bad"),
		},
		"blah",
	)
	if result {
		t.Fatal("Expected false, got true")
	}
}
