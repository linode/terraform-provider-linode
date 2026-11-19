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

const resourceName = "linode_nfs_filesystem_access_policy.foobar"

// minCapacityBytes is the minimum provisionable filesystem capacity (1 GiB).
const minCapacityBytes int64 = 1073741824

func TestAccResourceNFSFilesystemAccessPolicy_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(
					t, label, testRegion, minCapacityBytes,
					true, string(linodego.NFSSquashPolicyRootSquash),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("squash_policy"),
						knownvalue.StringExact(string(linodego.NFSSquashPolicyRootSquash))),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("inherit"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("linode_ids"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("protocols"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.Basic(
					t, label, testRegion, minCapacityBytes,
					false, string(linodego.NFSSquashPolicyAllSquash),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("squash_policy"),
						knownvalue.StringExact(string(linodego.NFSSquashPolicyAllSquash))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
