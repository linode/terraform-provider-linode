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

func TestAccLinodeNodeBalancerNodeBasic(t *testing.T) {
	// t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := fmt.Sprintf("tf_test_%s", acctest.RandString(10))
	config := testAccCheckLinodeNodeBalancerNodeBasic(nodeName)

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				ResourceName: resName,
				// ImportState:  true,
				// ImportStateVerify: true,
				Config: config,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", nodeName),
					resource.TestCheckResourceAttr(resName, "address", "192.168.200.1:80"),

					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttr(resName, "mode", "accept"),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
		},
	})
}

func TestAccLinodeNodeBalancerNodeUpdate(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerNodeBasic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", nodeName),
					resource.TestCheckResourceAttr(resName, "address", "192.168.200.1:80"),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeNodeBalancerNodeUpdates(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", nodeName)),
					resource.TestCheckResourceAttr(resName, "address", "192.168.200.1:8080"),
					resource.TestCheckResourceAttr(resName, "weight", "200"),
				),
			},
		},
	})
}

func testAccCheckLinodeNodeBalancerNodeExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_node" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		configID, err := strconv.Atoi(rs.Primary.Attributes["config_id"])

		_, err = client.GetNodeBalancerNode(nodebalancerID, configID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer Node %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeNodeBalancerNodeDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Failed to get Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_node" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		configID, err := strconv.Atoi(rs.Primary.Attributes["config_id"])

		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetNodeBalancerNode(nodebalancerID, configID, id)

		if err == nil {
			return fmt.Errorf("NodeBalancer Node with id %d still exists", id)
		}

		if apiErr, ok := err.(linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Failed to request NodeBalancer Node with id %d", id)
		}
	}

	return nil
}
func testAccCheckLinodeNodeBalancerNodeBasic(label string) string {
	return testAccCheckLinodeNodeBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "192.168.200.1:80"
	label = "%s"
	weight = 50
}
`, label)
}

func testAccCheckLinodeNodeBalancerNodeUpdates(label string) string {
	return testAccCheckLinodeNodeBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "192.168.200.1:8080"
	label = "%s_renamed"
	weight = 200
}

`, label)
}
