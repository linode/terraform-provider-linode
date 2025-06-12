//go:build integration || networkingipassignment

package networkingipassignment_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingipassignment/tmpl"
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

	resourceName := "linode_networking_ip_assignment.test"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNetworkingIPsAssignDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NetworkingIPsAssign(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.#"),
					func(*terraform.State) error {
						time.Sleep(30 * time.Second) // Add a delay to allow for API propagation
						return nil
					},
					checkNetworkingIPsAssignExists,
					resource.TestCheckResourceAttrSet(resourceName, "assignments.0.linode_id"),
					resource.TestCheckResourceAttrSet(resourceName, "assignments.0.address"),
				),
			},
			// Removed ImportState step as it's no longer supported
		},
	})
}

func checkNetworkingIPsAssignExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_networking_assign_ip" {
			continue
		}

		filter := fmt.Sprintf("{\"region\": \"%s\"}", rs.Primary.Attributes["region"])
		ips, err := client.ListIPAddresses(context.Background(), &linodego.ListOptions{Filter: filter})
		if err != nil {
			return fmt.Errorf("Error listing IP addresses: %s", err)
		}

		assignmentCount := 0
		for _, ip := range ips {
			if ip.LinodeID != 0 {
				assignmentCount++
			}
		}

		if assignmentCount == 0 {
			return fmt.Errorf("No IP assignments found")
		}
	}

	return nil
}

func checkNetworkingIPsAssignDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_networking_assign_ip" {
			continue
		}

		ipAddress := rs.Primary.ID // Assuming ID is the address of the IP
		_, err := client.GetIPAddress(context.Background(), ipAddress)
		if err == nil {
			return fmt.Errorf("Networking IPs Assign with id %s still exists", ipAddress)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Networking IPs Assign with id %s", ipAddress)
		}
	}

	return nil
}
