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

func init() {
	resource.AddTestSweepers("linode_nodebalancer", &resource.Sweeper{
		Name: "linode_nodebalancer",
		F:    testSweepLinodeNodeBalancer,
	})
}

func testSweepLinodeNodeBalancer(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	nodebalancers, err := client.ListNodeBalancers(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting instances: %s", err)
	}
	for _, nodebalancer := range nodebalancers {
		if nodebalancer.Label == nil || !shouldSweepAcceptanceTestResource(prefix, *nodebalancer.Label) {
			continue
		}
		err := client.DeleteNodeBalancer(context.Background(), nodebalancer.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %v during sweep: %s", nodebalancer.Label, err)
		}
	}

	return nil
}

func TestAccLinodeNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
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
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
				),
			},

			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				Config: testAccCheckLinodeNodeBalancerBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
				),
			},
			{
				Config: testAccCheckLinodeNodeBalancerUpdates(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", nodebalancerName)),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "0"),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.2667398925", "tf_test_2"),
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
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

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
	tags = ["tf_test"]
}
`, nodebalancer)
}

func testAccCheckLinodeNodeBalancerUpdates(nodebalancer string) string {
	return fmt.Sprintf(`
resource "linode_nodebalancer" "foobar" {
	label = "%s_r"
	region = "us-east"
	client_conn_throttle = 0
	tags = ["tf_test", "tf_test_2"]
}
`, nodebalancer)
}
