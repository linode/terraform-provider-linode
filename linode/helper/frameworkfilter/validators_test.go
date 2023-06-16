package frameworkfilter

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFilterableValidator_filter(t *testing.T) {
	v := testFilterConfig.validateFilterable(false)

	var d diag.Diagnostics

	response := validator.StringResponse{
		Diagnostics: d,
	}

	v.ValidateString(
		context.Background(),
		validator.StringRequest{
			ConfigValue: types.StringValue("fake"),
		},
		&response,
	)

	if !response.Diagnostics.HasError() {
		t.Fatal("expected an error; got none")
	}

	if response.Diagnostics[0].Summary() != "Non-Filterable Field" {
		t.Fatal("summary mismatch")
	}

	expectedError := "Field \"fake\" is not filterable.\nFilterable Fields: api_bar, api_foo, bar, foo"
	if response.Diagnostics[0].Detail() != expectedError {
		t.Fatal(cmp.Diff(response.Diagnostics[0].Detail(), expectedError))
	}
}

func TestFilterableValidator_order(t *testing.T) {
	v := testFilterConfig.validateFilterable(true)

	var d diag.Diagnostics

	response := validator.StringResponse{
		Diagnostics: d,
	}

	v.ValidateString(
		context.Background(),
		validator.StringRequest{
			ConfigValue: types.StringValue("foo"),
		},
		&response,
	)

	if !response.Diagnostics.HasError() {
		t.Fatal("expected an error; got none")
	}

	if response.Diagnostics[0].Summary() != "Unable to order by field" {
		t.Fatal("summary mismatch")
	}

	expectedError := "Field \"foo\" cannot be used in order_by as it is not API filterable.\n" +
		"API Filterable Fields: api_bar, api_foo"
	if response.Diagnostics[0].Detail() != expectedError {
		t.Fatal(cmp.Diff(response.Diagnostics[0].Detail(), expectedError))
	}
}

func TestFilterableValidator_success(t *testing.T) {
	v := testFilterConfig.validateFilterable(true)

	var d diag.Diagnostics

	response := validator.StringResponse{
		Diagnostics: d,
	}

	v.ValidateString(
		context.Background(),
		validator.StringRequest{
			ConfigValue: types.StringValue("api_foo"),
		},
		&response,
	)

	if response.Diagnostics.HasError() {
		t.Fatalf("expected no error; got %s", response.Diagnostics[0].Detail())
	}
}
