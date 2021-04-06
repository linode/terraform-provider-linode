package linode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testLKEClusterDataName = "data.linode_lke_cluster.test"

func TestAccDataSourceLinodeLKECluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeLKEClusterBasic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLKEClusterDataName, "label", clusterName),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "region", "us-central"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "k8s_version", "1.20"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "status", "ready"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "pools.#", "1"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "pools.0.type", "g6-standard-2"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "pools.0.count", "3"),
					resource.TestCheckResourceAttr(testLKEClusterDataName, "pools.0.nodes.#", "3"),
					resource.TestCheckResourceAttrSet(testLKEClusterDataName, "pools.0.id"),
					resource.TestCheckResourceAttrSet(testLKEClusterDataName, "kubeconfig"),
				),
			},
		},
	})
}

func testDataSourceLinodeLKEClusterBasic(clusterName string) string {
	return testAccCheckLinodeLKEClusterBasic(clusterName) + `
data "linode_lke_cluster" "test" {
	id = linode_lke_cluster.test.id
}`
}
