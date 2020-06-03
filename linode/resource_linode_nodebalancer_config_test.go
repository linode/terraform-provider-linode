package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func TestAccLinodeNodeBalancerConfig_basic(t *testing.T) {
	// t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	config := testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName)
	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckLinodeNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:       config,
				ResourceName: resName,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/"),

					resource.TestCheckResourceAttrSet(resName, "algorithm"),
					resource.TestCheckResourceAttrSet(resName, "stickiness"),
					resource.TestCheckResourceAttrSet(resName, "check_attempts"),
					resource.TestCheckResourceAttrSet(resName, "check_timeout"),
					resource.TestCheckResourceAttrSet(resName, "check_interval"),
					resource.TestCheckResourceAttrSet(resName, "check_passive"),
					resource.TestCheckResourceAttrSet(resName, "cipher_suite"),
					resource.TestCheckNoResourceAttr(resName, "ssl_common"),
					resource.TestCheckNoResourceAttr(resName, "ssl_ciphersuite"),
					resource.TestCheckResourceAttr(resName, "node_status.#", "1"),
					resource.TestCheckResourceAttr(resName, "node_status.0.up", "0"),
					resource.TestCheckResourceAttr(resName, "node_status.0.down", "0"),
					resource.TestCheckNoResourceAttr(resName, "ssl_cert"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateIDNodeBalancerConfig,
			},
		},
	})
}

func TestAccLinodeNodeBalancerConfig_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/"),
					resource.TestCheckResourceAttr(resName, "check_passive", "true"),

					resource.TestCheckResourceAttrSet(resName, "stickiness"),
					resource.TestCheckResourceAttrSet(resName, "check_attempts"),
					resource.TestCheckResourceAttrSet(resName, "check_timeout"),
				),
			},
			{
				Config: testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8088"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/foo"),
					resource.TestCheckResourceAttr(resName, "check_attempts", "3"),
					resource.TestCheckResourceAttr(resName, "check_timeout", "30"),
					resource.TestCheckResourceAttr(resName, "check_passive", "false"),

					resource.TestCheckResourceAttr(resName, "stickiness", string(linodego.StickinessHTTPCookie)),
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
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}

		_, err = client.GetNodeBalancerConfig(context.Background(), nodebalancerID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer Config %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeNodeBalancerConfigDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
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

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetNodeBalancerConfig(context.Background(), nodebalancerID, id)

		if err == nil {
			return fmt.Errorf("NodeBalancer Config with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting NodeBalancer Config with id %d", id)
		}
	}

	return nil
}

func testAccStateIDNodeBalancerConfig(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
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
		return fmt.Sprintf("%d,%d", nodebalancerID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_nodebalancer_config")
}

func testAccCheckLinodeNodeBalancerConfigBasic(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8080
	protocol = "http"
	check = "http"
	check_passive = true
	check_path = "/"
}
`
}

func testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8088
	protocol = "http"
	check = "http"
	check_path = "/foo"
	check_attempts = 3
	check_timeout = 30
	check_passive = false
	stickiness = "http_cookie"
	algorithm = "source"
}
`
}
