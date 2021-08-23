package balancernode_test

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
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestAccResourceNodeBalancerNode_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := acctest.RandomWithPrefix("tf_test")
	config := resourceConfigBasic(nodeName)

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		Providers:                 acceptance.TestAccProviders,
		CheckDestroy:              checkNodeBalancerNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(config, map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerNodeExists,
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
				ImportStateIdFunc: importResourceStateID,
			},
		},
	})
}

func TestAccResourceNodeBalancerNode_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_node.foonode"
	nodeName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkNodeBalancerNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.AccTestWithProvider(resourceConfigBasic(nodeName), map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", nodeName),
					resource.TestCheckResourceAttr(resName, "weight", "50"),
				),
			},
			{
				Config: acceptance.AccTestWithProvider(resourceConfigUpdates(nodeName), map[string]interface{}{
					acceptance.SkipInstanceReadyPollKey: true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerNodeExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", nodeName)),
					resource.TestCheckResourceAttr(resName, "weight", "200"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importResourceStateID,
			},
		},
	})
}

func checkNodeBalancerNodeExists(s *terraform.State) (err error) {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkNodeBalancerNodeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

func importResourceStateID(s *terraform.State) (string, error) {
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

func resourceBalancerConfigNetworking(instance string, pubkey string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	label = "%s"
	type = "g6-nanode-1"
	image = "linode/ubuntu18.04"
	region = "us-east"
	root_pass = "terraform-test"
	swap_size = 256
	private_ip = true
	authorized_keys = ["%s"]
	group = "tf_test"
}`, instance, pubkey)
}

func resourceBalancerConfig(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s"
	region = "us-east"
	client_conn_throttle = 20
	tags = ["tf_test"]
}
`, nodebalancer)
}

func resourceBalancerConfigBasic(nodebalancer string) string {
	return resourceBalancerConfig(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8080
	protocol = "HttP"
	check = "http"
	check_passive = true
	check_path = "/"
}
`
}

func resourceConfigBasic(label string) string {
	return resourceBalancerConfigNetworking(label, acceptance.PublicKeyMaterial) + resourceBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "${linode_instance.foobar.private_ip_address}:80"
	label = "%s"
	weight = 50
}
`, label)
}

func resourceConfigUpdates(label string) string {
	return resourceBalancerConfigNetworking(label, acceptance.PublicKeyMaterial) + resourceBalancerConfigBasic(label) + fmt.Sprintf(`
resource "linode_nodebalancer_node" "foonode" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	config_id = "${linode_nodebalancer_config.foofig.id}"
	address = "${linode_instance.foobar.private_ip_address}:8080"
	label = "%s_r"
	weight = 200
}

`, label)
}
