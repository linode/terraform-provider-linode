package firewall

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-set/v3"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var firewallProtocolKeywords = set.From([]string{"ALL", "TCP", "UDP", "ICMP", "IPENCAP"})

var _ validator.String = firewallProtocolValidator{}

type firewallProtocolValidator struct{}

func (v firewallProtocolValidator) Description(ctx context.Context) string {
	return "value must be ALL, TCP, UDP, ICMP, IPENCAP, or a protocol number from 0 to 255"
}

func (v firewallProtocolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v firewallProtocolValidator) ValidateString(
	ctx context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	protocol := req.ConfigValue.ValueString()
	if isValidFirewallProtocol(protocol) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Firewall Rule Protocol",
		fmt.Sprintf(
			"Expected ALL, TCP, UDP, ICMP, IPENCAP, or a protocol number from 0 to 255 without leading zeros, got %q.",
			protocol,
		),
	)
}

func isValidFirewallProtocol(protocol string) bool {
	if firewallProtocolKeywords.Contains(protocol) {
		return true
	}

	value, err := strconv.Atoi(protocol)
	if err != nil {
		return false
	}

	return value >= 0 && value <= 255 && strconv.Itoa(value) == protocol
}
