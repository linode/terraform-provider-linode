//go:build unit

package helper_test

import (
	"strings"
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func TestSDKv2ValidateIPv4Range(t *testing.T) {
	d := helper.SDKv2ValidateIPv4Range("192.168.0.1", nil)
	if d == nil || !d.HasError() {
		t.Fatal("Expected error, got none")
	}

	if !strings.Contains(d[0].Summary, "Invalid IPv4 CIDR range: 192.168.0.1") {
		t.Fatalf("Error does not match expected error: %s", d[0].Summary)
	}

	d = helper.SDKv2ValidateIPv4Range("::0/0", nil)
	if d == nil || !d.HasError() {
		t.Fatal("Expected error, got none")
	}

	if !strings.Contains(d[0].Summary, "Expected IPv4 address, got IPv6") {
		t.Fatal("Error does not match expected error")
	}

	d = helper.SDKv2ValidateIPv4Range("192.168.0.0/24", nil)
	if d != nil && d.HasError() {
		t.Fatal("Expected none, got error")
	}
}

func TestSDKv2ValidateIPv6Range(t *testing.T) {
	d := helper.SDKv2ValidateIPv6Range("::0", nil)
	if d == nil || !d.HasError() {
		t.Fatal("Expected error, got none")
	}

	if !strings.Contains(d[0].Summary, "Invalid IPv6 CIDR range: ::0") {
		t.Fatalf("Error does not match expected error: %s", d[0].Summary)
	}

	d = helper.SDKv2ValidateIPv6Range("192.168.0.1/24", nil)
	if d == nil || !d.HasError() {
		t.Fatal("Expected error, got none")
	}

	if !strings.Contains(d[0].Summary, "Expected IPv6 address, got IPv4") {
		t.Fatal("Error does not match expected error")
	}

	d = helper.SDKv2ValidateIPv6Range("::0/0", nil)
	if d != nil && d.HasError() {
		t.Fatal("Expected none, got error")
	}
}
