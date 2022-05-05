package helper_test

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestCompareIPv6Ranges(t *testing.T) {
	ips := []string{
		"1111:1111::1111:1111:1111:f88/128",
		"1111:1111::1111:1111:1111:0f88/128",
	}

	result, err := helper.CompareIPv6Ranges(ips[0], ips[1])
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Fatalf("ranges are reported as different despite being equal")
	}

	ips[1] = "1111:1111::1111:1211:1111:0f88/126"
	result, err = helper.CompareIPv6Ranges(ips[0], ips[1])
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Fatalf("ranges are reported as equal despite being different")
	}
}
