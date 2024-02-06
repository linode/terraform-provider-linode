//go:build integration

package nodepool_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nodepool"
	"log"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/nodepool/tmpl"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_nodepool", &resource.Sweeper{
		Name: "linode_nodepool",
		F:    sweep,
	})
}

func GetTestClient() (*linodego.Client, error) {
	token := os.Getenv("LINODE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("LINODE_TOKEN must be set for acceptance tests")
	}

	apiVersion := os.Getenv("LINODE_API_VERSION")
	if apiVersion == "" {
		apiVersion = "v4beta"
	}

	config := &helper.Config{
		AccessToken: token,
		APIVersion:  apiVersion,
		APIURL:      os.Getenv("LINODE_URL"),
	}

	client, err := config.Client(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func sweep(prefix string) error {
	client, err := GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}
	clusterID, err := strconv.Atoi(os.Getenv("LINODE_TEST_CLUSTER_ID"))

	pools, err := client.ListLKENodePools(context.Background(), clusterID, nil)
	if err != nil {
		return fmt.Errorf("Error getting node pools: %s", err)
	}
	for _, nodepool := range pools {
		if containsTagWithPrefix(nodepool, "pool_test_") {
			log.Printf("[DEBUG] Found a leaked node pool, clusterID: %d, poolID: %d. Deleting", clusterID, nodepool.ID)
			err := client.DeleteLKENodePool(context.Background(), clusterID, nodepool.ID)
			if err != nil {
				return fmt.Errorf("Error destroying nodepool %v during sweep: %s", nodepool.ID, err)
			}
		}
	}
	return nil
}

func TestAccResourceNodePool_basic(t *testing.T) {
	t.Parallel()

	clusterID := os.Getenv("LINODE_TEST_CLUSTER_ID")
	tag := acctest.RandomWithPrefix("pool_test_")
	resName := "linode_nodepool.foobar"

	poolTestPreCheck := func() {
		acceptance.PreCheck(t)
		if v := os.Getenv("LINODE_TEST_CLUSTER_ID"); v == "" {
			t.Fatal("LINODE_TEST_CLUSTER_ID must be set for acceptance tests")
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 poolTestPreCheck,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, clusterID, tag),
				Check: resource.ComposeTestCheckFunc(
					checkNodePoolExists,
					resource.TestCheckResourceAttr(resName, "node_count", "2"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-4"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", tag),
					resource.TestCheckResourceAttr(resName, "autoscaler.min", "2"),
					resource.TestCheckResourceAttr(resName, "autoscaler.max", "3"),
				),
			},
		},
	})
}

func checkNodePoolExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodepool" {
			continue
		}

		clusterID, poolID, err := nodepool.ParseNodePoolID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetLKENodePool(context.Background(), clusterID, poolID)
		if err != nil {
			return fmt.Errorf("Error retrieving state of node pool %v: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

func checkNodePoolDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodepool" {
			continue
		}

		clusterID, poolID, err := nodepool.ParseNodePoolID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetLKENodePool(context.Background(), clusterID, poolID)

		if err == nil {
			return fmt.Errorf("Node Pool with id %d still exists in cluster %d", poolID, clusterID)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Node Pool with id %d in cluster %d", poolID, clusterID)
		}
	}

	return nil
}

func containsTagWithPrefix(pool linodego.LKENodePool, prefix string) bool {
	for _, tag := range pool.Tags {
		if strings.HasPrefix(tag, prefix) {
			return true
		}
	}
	return false
}
