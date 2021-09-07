package lke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/lke/tmpl"
)

const dataSourceClusterName = "data.linode_lke_cluster.test"

func TestAccDataSourceLinodeLKECluster_basic(t *testing.T) {
	t.Parallel()

	clusterName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckLKEClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
					resource.TestCheckResourceAttr(dataSourceClusterName, "region", "us-central"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", "1.20"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-2"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "3"),
					resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.nodes.#", "3"),
					resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),
				),
			},
		},
	})
}
