package helper

import (
	"context"
	"fmt"
	"net"
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
	return fmt.Sprintf(
		"value must meet RFC3339 standard and in format of '%s'.",
		TIME_FORMAT,
	)
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

type IPStringValidator struct{}

func (ipv IPStringValidator) Description(ctx context.Context) string {
	return ipv.MarkdownDescription(ctx)
}

func (ipv IPStringValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Value must be a valid IPv4 or IPv6 IP Address.")
}

func (ipv IPStringValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	v := request.ConfigValue.ValueString()

	ip := net.ParseIP(v)

	if ip == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			ipv.Description(ctx),
			v,
		))
	}
}

func NewIPStringValidator() IPStringValidator {
	return IPStringValidator{}
}

type StringLengthValidator struct {
	Minimum int
	Maximum int
}

func (slv StringLengthValidator) Description(ctx context.Context) string {
	return slv.MarkdownDescription(ctx)
}

func (slv StringLengthValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Value must have length within range [%d,%d] (inclusive).", slv.Minimum, slv.Maximum)
}

func (slv StringLengthValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	v := request.ConfigValue.ValueString()

	if len(v) < slv.Minimum || len(v) > slv.Maximum {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			slv.Description(ctx),
			v,
		))
	}
}

func NewStringLengthValidator(min, max int) StringLengthValidator {
	return StringLengthValidator{Minimum: min, Maximum: max}
}
