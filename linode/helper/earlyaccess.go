package helper

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// AttemptWarnEarlyAccessSDKv2 will raise a warning if an SDKv2
// early access resource is being used without the v4beta API version.
func AttemptWarnEarlyAccessSDKv2(providerMeta *ProviderMeta) {
	if providerMeta.Config.APIVersion == "v4beta" {
		return
	}

	log.Printf(
		"[WARN] This resource is in early access but the provider "+
			"API version is set to \"%s\" (expected \"v4beta\").",
		providerMeta.Config.APIVersion,
	)
}

// AttemptWarnEarlyAccessFramework will raise a warning if a Framework
// early access resource is being used without the v4beta API version.
func AttemptWarnEarlyAccessFramework(config *FrameworkProviderModel) diag.Diagnostics {
	var d diag.Diagnostics

	if config.APIVersion.ValueString() != "v4beta" {
		d.AddWarning(
			"Non-Beta Target API Version",
			fmt.Sprintf(
				"This resource is in early access but the provider "+
					"API version is set to \"%s\" (expected \"v4beta\").",
				config.APIVersion.ValueString(),
			),
		)
	}

	return d
}
