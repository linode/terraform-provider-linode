package frameworkfilter

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type filterNameValidator struct {
	FilterConfig Config
	IsOrderBy    bool
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
	if req.ConfigValue.IsNull() {
		return
	}

	fieldName := req.ConfigValue.ValueString()

	config, ok := v.FilterConfig[fieldName]

	if !ok {
		// Aggregate filterable attributes
		var filterableAttributes []string

		for k := range v.FilterConfig {
			filterableAttributes = append(filterableAttributes, k)
		}

		sort.Strings(filterableAttributes)

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Non-Filterable Field",
			fmt.Sprintf(
				"Field \"%s\" is not filterable.\nFilterable Fields: %s",
				fieldName,
				strings.Join(filterableAttributes, ", "),
			),
		)

		return
	}

	if v.IsOrderBy && !config.APIFilterable && !config.AllowOrderOverride {
		// Aggregate filterable attributes
		var filterableAttributes []string

		for k, v := range v.FilterConfig {
			if !v.APIFilterable && !v.AllowOrderOverride {
				continue
			}

			filterableAttributes = append(filterableAttributes, k)
		}

		sort.Strings(filterableAttributes)

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Unable to order by field",
			fmt.Sprintf(
				"Field \"%s\" cannot be used in order_by as it is not API filterable.\n"+
					"API Filterable Fields: %s",
				fieldName,
				strings.Join(filterableAttributes, ", "),
			),
		)
	}
}

func (f Config) validateFilterable(IsOrderBy bool) filterNameValidator {
	return filterNameValidator{
		FilterConfig: f,
		IsOrderBy:    IsOrderBy,
	}
}
