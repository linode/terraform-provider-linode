//go:build unit

package helper_test

import (
	"testing"
	"context"

	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/sshkey"
	"github.com/linode/terraform-provider-linode/linode/nb"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRegexSuccess_sshKey(t *testing.T) {
	pattern := sshkey.SSHKeyLabelRegex

	regExp := helper.StringToRegex(pattern)

	testValidStrings := []string{
		"valid_String123",
		"valid__string",
		"valid_string-WITH_dash",
	}

	testInvalidStrings := []string{
		"*InvalidString",
		"AnotherInvalid!",
		"Not..Invalid",
		"&notValid",
		"#Nope",
		"(NotValid)",
	}

	for _, str := range testValidStrings {
		if !regExp.MatchString(str) {
			t.Fatal("Should match regex")
		}
	}

	for _, str := range testInvalidStrings {
		if regExp.MatchString(str) {
			t.Fatal("Should not match regex")
		}
	}
}

func TestCheckSuccess_nbLabel(t *testing.T) {
	pattern := nb.NBLabelRegex

	regExp := helper.StringToRegex(pattern)

	testValidStrings := []string{
		"valid_String123",
		"valid__string",
		"valid_string-WITH_dash",
	}

	testInvalidStrings := []string{
		"*InvalidString",
		"AnotherInvalid!",
		"Not..Invalid",
		"&notValid",
		"#Nope",
		"(NotValid)",
	}

	for _, str := range testValidStrings {
		if !regExp.MatchString(str) {
			t.Fatal("Should match regex")
		}
	}

	for _, str := range testInvalidStrings {
		if regExp.MatchString(str) {
			t.Fatal("Should not match regex")
		}
	}
}

func TestRegexValidator_success(t *testing.T) {
	v := helper.RegexMatches(sshkey.SSHKeyLabelRegex, sshkey.SSHKeyLabelErrorMessage)

	var d diag.Diagnostics

	response := validator.StringResponse{
		Diagnostics: d,
	}

	v.ValidateString(
		context.Background(),
		validator.StringRequest{
			ConfigValue: types.StringValue("valid_String123"),
		},
		&response,
	)

	if response.Diagnostics.HasError() {
		t.Fatalf("expected no error; got %s", response.Diagnostics[0].Detail())
	}
}
