package helper

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

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
