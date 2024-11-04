package networkipassignment_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/networkipassignment/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceNetworkingIPsAssign(t *testing.T) {
	t.Parallel()

	resourceName := "linode_networking_ips_assign.test"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNetworkingIPsAssignDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPsAssign(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNetworkingIPsAssignExists,
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "assignments.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.0.address"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.1.linode_id"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.1.address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkNetworkingIPsAssignExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_networking_ips_assign" {
			continue
		}

		_, err := client.GetIPAddress(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Networking IPs Assign %s: %s", rs.Primary.Attributes["id"], err)
		}
	}

	return nil
}

func checkNetworkingIPsAssignDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_networking_ips_assign" {
			continue
		}

		_, err := client.GetIPAddress(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Networking IPs Assign with id %s still exists", rs.Primary.ID)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Networking IPs Assign with id %s", rs.Primary.ID)
		}
	}

	return nil
}
