package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_lke_cluster", &resource.Sweeper{
		Name: "linode_lke_cluster",
		F:    testSweepLinodeLKE,
	})
}

func testSweepLinodeLKE(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	clustersLKE, err := client.ListLKEClusters(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting templates: %s", err)
	}
	for _, lke := range clustersLKE {
		if !shouldSweepAcceptanceTestResource(prefix, lke.Label) {
			continue
		}
		err := client.DeleteLKECluster(context.Background(), lke.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", lke.Label, err)
		}
	}

	return nil
}

func TestAccLinodeLKE_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_cluster.foobar"
	lkeName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterConfigBasic(lkeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLKEClusterExists,
					resource.TestCheckResourceAttr(resName, "label", lkeName),
					resource.TestCheckResourceAttr(resName, "region", "us-central"),
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

func TestAccLinodeLKE_update(t *testing.T) {
	t.Parallel()

	resName := "linode_lke_cluster.foobar"
	lkeName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeLKEClusterConfigBasic(lkeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLKEClusterExists,
					resource.TestCheckResourceAttr(resName, "label", lkeName),
				),
			},
			{
				Config: testAccCheckLinodeLKEClusterUpdates(lkeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeLKEClusterExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_u", lkeName)),
				),
			},
		},
	})
}

func testAccCheckLinodeLKEClusterDelete(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lke_cluster" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetLKECluster(context.Background(), id)
		if err == nil {
			return fmt.Errorf("Linode LKE with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode LKE with id %d", id)
		}
	}
	return nil
}

func testAccCheckLinodeLKEClusterExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lke_cluster" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetLKECluster(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of LKE %s: %s", rs.Primary.Attributes["label"], err)
		}
	}
	return nil
}

// TODO(sgmac): test passes, destroy leaves linode nodes behind even though
// the cluster is destroyed.
func testAccCheckLinodeLKEClusterConfigBasic(label string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "foobar" {
	label = "%s"
	region = "us-central"
	version = "1.16"
	node_pools = [
		{ "count" = 3, "type" = "g6-standard-1"}
	]
}`, label)
}

func testAccCheckLinodeLKEClusterUpdates(label string) string {
	return fmt.Sprintf(`
resource "linode_lke_cluster" "foobar" {
	label = "%s_u"
	region = "us-central"
	version = "1.16"
	node_pools = [
		{ "count" = 3, "type" = "g6-standard-1"}
	]
}`, label)
}
