//go:build integration || nfsfilesystemaccesspolicy

package nfsfilesystemaccesspolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystemaccesspolicy/tmpl"
)

const dataSourceName = "data.linode_nfs_filesystem_access_policy.foobar"

func TestAccDataSourceNFSFilesystemAccessPolicy_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// A filesystem's access policy is unconfigured until it is first updated.
				Config: tmpl.DataUnconfigured(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("configured"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("enabled"), knownvalue.Null()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("status"), knownvalue.Null()),
				},
			},
			{
				Config: tmpl.DataConfigured(
					t, label, testRegion, minCapacityBytes,
					true, true, string(linodego.NFSSquashPolicyRootSquash),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("configured"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("inherit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("squash_policy"),
						knownvalue.StringExact(string(linodego.NFSSquashPolicyRootSquash))),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("protocols"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
		},
	})
}
