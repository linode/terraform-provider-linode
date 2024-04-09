//go:build integration || linode

package linode_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

func TestCreatingFrameworkProvider(t *testing.T) {
	_ = linode.CreateFrameworkProvider("test")
}

func TestAccFrameworkProvider_AlternativeEndpoint(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: alternativeAPIURLTemplate(
					"https://api.linode.com",
					"v4_cooler_version",
				),
			},
		},
	},
	)
}

func alternativeAPIURLTemplate(
	url string,
	apiVersion string,
) string {
	return fmt.Sprintf(`
provider "linode" {
  url = "%s"
  api_version = "%s"
}
`, url, apiVersion) // lintignore:AT004
}
