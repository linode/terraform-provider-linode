//go:build integration || nfsspaceaccesspolicy

package nfsspaceaccesspolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsspaceaccesspolicy/tmpl"
)

const dataSourceName = "data.linode_nfs_space_access_policy.foobar"

func TestAccDataSourceNFSSpaceAccessPolicy_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion, minCapacityBytes, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("mtls_mode"),
						knownvalue.StringExact(string(linodego.NFSMTLSModeDisabled))),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("vpc_acl"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
		},
	})
}
