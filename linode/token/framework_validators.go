package token

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type DateTimeStringValidator struct {
	Format string
}

func (dtv DateTimeStringValidator) Description(ctx context.Context) string {
	return dtv.MarkdownDescription(ctx)
}

func (dtv DateTimeStringValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("value must meet ISO 8601 standard format, e.g., '%s'.", time.RFC3339)
}

func (dtv DateTimeStringValidator) ValidateString(
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

func NewDateTimeStringValidator(format string) DateTimeStringValidator {
	return DateTimeStringValidator{Format: format}
}
