package helper

import (
	"net"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func SDKv2ValidateIPv4Range(i any, path cty.Path) diag.Diagnostics {
	ip, _, err := net.ParseCIDR(i.(string))
	if err != nil {
		return diag.Errorf("Invalid IPv4 CIDR range: %s", i)
	}

	if ip.To4() == nil {
		return diag.Errorf("Expected IPv4 address, got IPv6")
	}

	return nil
}

func SDKv2ValidateIPv6Range(i any, path cty.Path) diag.Diagnostics {
	ip, _, err := net.ParseCIDR(i.(string))
	if err != nil {
		return diag.Errorf("Invalid IPv6 CIDR range: %s", i)
	}

	if ip.To4() != nil {
		return diag.Errorf("Expected IPv6 address, got IPv4")
	}

	return nil
}
