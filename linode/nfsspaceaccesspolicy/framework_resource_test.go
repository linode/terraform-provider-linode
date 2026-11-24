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

const resourceName = "linode_nfs_space_access_policy.foobar"

// minCapacityBytes is the minimum provisionable filesystem capacity (1 GiB).
const minCapacityBytes int64 = 1073741824

func TestAccResourceNFSSpaceAccessPolicy_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, minCapacityBytes, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("space_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("mtls_mode"),
						knownvalue.StringExact(string(linodego.NFSMTLSModeDisabled))),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("vpcs"), knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.Basic(t, label, testRegion, minCapacityBytes, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("enabled"), knownvalue.Bool(false)),
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

func TestAccResourceNFSSpaceAccessPolicy_vpcs(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

	vpcs := tfjsonpath.New("vpcs")
	vpcACL := tfjsonpath.New("vpc_acl")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPC(t, label, testRegion, minCapacityBytes, "vpc"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, vpcs, knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, vpcACL, knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName,
						vpcACL.AtSliceIndex(0).AtMapKey("range"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						vpcACL.AtSliceIndex(0).AtMapKey("subnets"), knownvalue.ListSizeExact(0)),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tmpl.VPC(t, label, testRegion, minCapacityBytes, "subnet"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, vpcs, knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName,
						vpcACL.AtSliceIndex(0).AtMapKey("subnets"), knownvalue.ListSizeExact(1)),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tmpl.VPC(t, label, testRegion, minCapacityBytes, "empty"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, vpcs, knownvalue.SetSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName, vpcACL, knownvalue.ListSizeExact(0)),
				},
			},
		},
	})
}
