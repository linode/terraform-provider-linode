package helper

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// SDKv2ValidateFieldRequiresAPIVersion is an SDKv2 CustomizeDiffFunc
// that ensures the given fields aren't specified if the configuration API
// version does not match the required version.
//
// NOTE: This needs to be implemented as a CustomizeDiffFunc because it
// requires access to the provider config.
func SDKv2ValidateFieldRequiresAPIVersion(
	requiredAPIVersion string,
	fieldPaths ...string,
) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		providerMeta := meta.(*ProviderMeta)

		if strings.EqualFold(providerMeta.Config.APIVersion, requiredAPIVersion) {
			return nil
		}

		for _, fieldPath := range fieldPaths {
			_, newValue := diff.GetChange(fieldPath)
			newValueString := newValue.(string)

			if newValueString == "" {
				continue
			}

			return fmt.Errorf(
				"%s: The api_version provider argument must be set to '%s' to use this field.",
				fieldPath,
				requiredAPIVersion,
			)
		}

		return nil
	}
}
