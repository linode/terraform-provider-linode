//go:build unit

package helper_test

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestCheckRegex(t *testing.T) {
	pattern := "^[a-zA-Z0-9]([-_.]?[a-zA-Z0-9]+)*[a-zA-Z0-9]$"

	regExp := helper.StringToRegex(pattern)

	testValidStrings := []string{
		"valid_String123",
		"valid_string.with_period",
		"valid_string-with_dash",
	}

	testInvalidStrings := []string{
		"_InvalidString",
		"AnotherInvalid_",
		"Not..Invalid",
		"no_double--dash",
		"no_double__underscore",
		"!NotValid",
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
