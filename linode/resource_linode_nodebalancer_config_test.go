package linode

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLinodeNodeBalancerConfigBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foobar"
	nodebalancerName := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),

					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeNodeBalancerConfigUpdate(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "algorithm", "source"),
				),
			},
		},
	})
}

func testAccCheckLinodeNodeBalancerConfigExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])

		_, err = client.GetNodeBalancerConfig(nodebalancerID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer Config %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeNodeBalancerConfigDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Failed to get Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])

		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetNodeBalancerConfig(nodebalancerID, id)

		if err == nil {
			return fmt.Errorf("NodeBalancer Config with id %d still exists", id)
		}

		if apiErr, ok := err.(linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Failed to request NodeBalancer with id %d", id)
		}
	}

	return nil
}
func testAccCheckLinodeNodeBalancerConfigBasic(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s"
	region = "us-east"
	client_conn_throttle = 20
}

resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
}
`, nodebalancer)
}

func testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s"
	region = "us-east"
	client_conn_throttle = 20
}
		
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	algorithm = "source"
}

`, nodebalancer)
}
