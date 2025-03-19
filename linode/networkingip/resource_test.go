//go:build integration || networkingip

package networkingip_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip/tmpl"
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceNetworkingIP_ephemeral(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resourceName := "linode_networking_ip.reserved_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPReservedAssigned(t, label, testRegion, 0, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "reserved", "false"),
					resource.TestCheckResourceAttr(resourceName, "public", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway"),
					resource.TestCheckResourceAttrPair(resourceName, "linode_id", "linode_instance.test.0", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "rdns"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_mask"),
					resource.TestCheckNoResourceAttr(resourceName, "vpc_nat_1_1"),
				),
			},
		},
	})
}

func TestAccResourceNetworkingIP_reserved(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resourceName := "linode_networking_ip.reserved_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPReservedAssigned(t, label, testRegion, 0, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "reserved", "true"),
					resource.TestCheckResourceAttr(resourceName, "public", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway"),
					resource.TestCheckResourceAttrPair(resourceName, "linode_id", "linode_instance.test.0", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "rdns"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_mask"),
					resource.TestCheckNoResourceAttr(resourceName, "vpc_nat_1_1"),
				),
			},
		},
	})
}

func TestAccResourceNetworkingIP_reservedEphemeralReassignment(t *testing.T) {
	t.Parallel()

	resName := "linode_networking_ip.reserved_ip"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			// Create an assigned reserved IP
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resName, "linode_id",
						"linode_instance.test.0", "id",
					),
					resource.TestCheckResourceAttr(resName, "reserved", "true"),
				),
			},
			// Make the IP ephemeral
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					0,
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						resName, "linode_id",
						"linode_instance.test.0", "id",
					),
					resource.TestCheckResourceAttr(resName, "reserved", "false"),
				),
			},
			// Attempt to reassign the ephemeral IP; expect RequiresReplace
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					1,
					false,
				),
				PlanOnly: true,
				// This plan is expected to trigger a RequiresReplace
				ExpectNonEmptyPlan: true,
			},
			// Convert back to a reserved IP, assign to second instance
			{
				Config: tmpl.NetworkingIPReservedAssigned(
					t,
					linodeLabel,
					testRegion,
					1,
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "reserved", "true"),
					resource.TestCheckResourceAttrPair(
						resName, "linode_id",
						"linode_instance.test.1", "id",
					),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"wait_for_available"},
			},
		},
	})
}
