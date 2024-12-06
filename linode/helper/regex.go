package helper

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringToRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func RegexMatches(pattern string, errorMessage string) validator.String {
	patternRegEx := StringToRegex(pattern)
	return stringvalidator.RegexMatches(patternRegEx, errorMessage)
}
