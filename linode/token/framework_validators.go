package token

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"time"
)

type dateTimeStringValidator struct {
	Format string
}

func (dtv dateTimeStringValidator) Description(ctx context.Context) string {
	return dtv.MarkdownDescription(ctx)
}

func (dtv dateTimeStringValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("value must meet ISO 8601 standard format, e.g., '%s'.", time.RFC3339)
}

func (dtv dateTimeStringValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if dtv.Format == "" {
		dtv.Format = time.RFC3339
	}
	v := request.ConfigValue.ValueString()
	if _, err := time.Parse(time.RFC3339, v); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			dtv.Description(ctx),
			v,
		))
	}
}

func DateTimeStringValidator(format string) validator.String {
	return dateTimeStringValidator{Format: format}
}
