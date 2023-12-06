package helper

import (
	"net"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

func SDKv2ObjectCannedACLValidator(i any, p cty.Path) diag.Diagnostics {
	aclValues, err := StringAliasSliceToStringSlice[s3types.ObjectCannedACL](
		s3types.ObjectCannedACLPrivate.Values(), // this return all acl values, not just private
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return validation.ToDiagFunc(validation.StringInSlice(aclValues, true))(i, p)
}
