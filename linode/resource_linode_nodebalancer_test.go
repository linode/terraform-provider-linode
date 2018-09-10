package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),

					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "transfer.%", "3"),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeNodeBalancer_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerUpdates(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", nodebalancerName)),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "0"),
				),
			},
		},
	})
}

func testAccCheckLinodeNodeBalancerExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		_, err = client.GetNodeBalancer(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeNodeBalancerDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetNodeBalancer(context.Background(), id)

		if err == nil {
			return fmt.Errorf("NodeBalancer with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting NodeBalancer with id %d", id)
		}
	}

	return nil
}
func testAccCheckLinodeNodeBalancerBasic(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s"
	region = "us-east"
	client_conn_throttle = 20
}
`, nodebalancer)
}

func testAccCheckLinodeNodeBalancerUpdates(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s_r"
	region = "us-east"
	client_conn_throttle = 0
}
`, nodebalancer)
}
