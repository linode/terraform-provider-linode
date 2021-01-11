package linode

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeNodeBalancerNode_basic(t *testing.T) {
	// t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := acctest.RandomWithPrefix("tf_test")
	config := testAccCheckLinodeNodeBalancerNodeBasic(nodeName)

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNode,
					resource.TestCheckResourceAttr(resName, "label", nodeName),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttr(resName, "mode", "accept"),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateIDNodeBalancerNode,
			},
		},
	})
}

func TestAccLinodeNodeBalancerNode_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeNodeBalancerNodeBasic(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNode,
					resource.TestCheckResourceAttr(resName, "label", nodeName),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
			{
				Config: testAccCheckLinodeNodeBalancerNodeUpdates(nodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerNode,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", nodeName)),
					resource.TestCheckResourceAttr(resName, "weight", "200"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateIDNodeBalancerNode,
			},
		},
	})
}

func testAccCheckLinodeNodeBalancerNode(s *terraform.State) (err error) {
	client := testAccProvider.Meta().(*ProviderMeta).Client

	var linodeID, nodebalancerID, nodeID, configID int
	var expectedNodePort string

	// find Linode instance ID
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance" {
			continue
		}

		linodeID, err = strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
	}

	// find NodeBalancer Node ID
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_node" {
			continue
		}

		nodeID, err = strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		nodebalancerID, err = strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}

		configID, err = strconv.Atoi(rs.Primary.Attributes["config_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["config_id"])
		}

		expectedNodePort = strings.Split(rs.Primary.Attributes["address"], ":")[1]
	}

	instanceNetwork, err := client.GetInstanceIPAddresses(context.Background(), linodeID)
	if err != nil {
		return fmt.Errorf("failed to get IPs for instance %d: %s", linodeID, err)
	}

	node, err := client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, nodeID)
	if err != nil {
		return fmt.Errorf("Error retrieving state of NodeBalancer Node %d: %s", nodeID, err)
	}

	privateIP := instanceNetwork.IPv4.Private[0].Address

	nodeAddrComps := strings.Split(node.Address, ":")
	nodeHost, nodePort := nodeAddrComps[0], nodeAddrComps[1]

	if nodeHost != privateIP {
		return fmt.Errorf("expected node to have host '%s'; got '%s'", privateIP, node.Address)
	}

	if nodePort != expectedNodePort {
		return fmt.Errorf("expected node to have port '%s'; got '%s'", expectedNodePort, nodePort)
	}
	return
}

func testAccCheckLinodeNodeBalancerNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_node" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}

		configID, err := strconv.Atoi(rs.Primary.Attributes["config_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["config_id"])
		}

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetNodeBalancerNode(context.Background(), nodebalancerID, configID, id)

		if err == nil {
			return fmt.Errorf("NodeBalancer Node with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting NodeBalancer Node with id %d", id)
		}
	}

	return nil
}

func testAccStateIDNodeBalancerNode(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_node" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing nodebalancer_id %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}
		configID, err := strconv.Atoi(rs.Primary.Attributes["config_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing config_id %v to int", rs.Primary.Attributes["config_id"])
		}
		return fmt.Sprintf("%d,%d,%d", nodebalancerID, configID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_nodebalancer_config")
}

func testAccCheckLinodeNodeBalancerNodeBasic(label string) string {
	return testAccCheckLinodeInstanceConfigPrivateNetworking(label, publicKeyMaterial) + testAccCheckLinodeNodeBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "${linode_instance.foobar.private_ip_address}:80"
	label = "%s"
	weight = 50
}
`, label)
}

func testAccCheckLinodeNodeBalancerNodeUpdates(label string) string {
	return testAccCheckLinodeInstanceConfigPrivateNetworking(label, publicKeyMaterial) + testAccCheckLinodeNodeBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "${linode_instance.foobar.private_ip_address}:8080"
	label = "%s_r"
	weight = 200
}

`, label)
}
