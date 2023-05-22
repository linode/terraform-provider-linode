package frameworkfilter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type filterNameValidator struct {
	FilterConfig Config
}

func (v filterNameValidator) Description(ctx context.Context) string {
	return "validate that the provided field is filterable"
}

func (v filterNameValidator) MarkdownDescription(ctx context.Context) string {
	return "validate that the provided field is filterable"
}

func (v filterNameValidator) ValidateString(
	ctx context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	fieldName := req.ConfigValue.ValueString()

	if _, ok := v.FilterConfig[fieldName]; !ok {
		// Aggregate filterable attributes
		filterableAttributes := make([]string, len(v.FilterConfig))
		i := 0

		for k, _ := range v.FilterConfig {
			filterableAttributes[i] = k
			i++
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Non-Filterable Field",
			fmt.Sprintf(
				"Field \"%s\" is not filterable.\nFilterable fields: %s",
				fieldName,
				strings.Join(filterableAttributes, ",")),
		)
	}
}

func (f Config) validateFilterable() filterNameValidator {
	return filterNameValidator{
		FilterConfig: f,
	}
}
