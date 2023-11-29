package helper

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
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

type regexValidator struct {
	regexp  *regexp.Regexp
}

func (v regexValidator) Description(ctx context.Context) string {
	return "validate that the provided field conforms to the regEx"
}

func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return "validate that the provided field conforms to the regEx"
}

func (v regexValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// MatchesRegex checks that the String held in the attribute
// matches the given RegEx.
func MatchesRegex(pattern string) regexValidator {
	return regexValidator{
		regexp: StringToRegex(pattern),
	}
}
