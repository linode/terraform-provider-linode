package helper

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func StringToRegex(pattern string) (regExp *regexp.Regexp) {
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	return regExp
}

func RegexMatches(pattern string, errorMessage string) validator.String {
	patternRegEx := StringToRegex(pattern)
	return stringvalidator.RegexMatches(patternRegEx, errorMessage)
}
