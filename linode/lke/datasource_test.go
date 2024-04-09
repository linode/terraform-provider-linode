//go:build integration || lke

package lke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lke/tmpl"
)

const dataSourceClusterName = "data.linode_lke_cluster.test"

func TestAccDataSourceLKECluster_basic(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 2, func(tRetry *acceptance.TRetry) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, clusterName, k8sVersionLatest, testRegion),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
						resource.TestCheckResourceAttr(dataSourceClusterName, "region", testRegion),
						resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", k8sVersionLatest),
						resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-2"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "3"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.nodes.#", "3"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.autoscaler.#", "0"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.high_availability", "false"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.id"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "dashboard_url"),
					),
				},
			},
		})
	})
}

func TestAccDataSourceLKECluster_autoscaler(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 2, func(tRetry *acceptance.TRetry) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataAutoscaler(t, clusterName, k8sVersionLatest, testRegion),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
						resource.TestCheckResourceAttr(dataSourceClusterName, "region", testRegion),
						resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", k8sVersionLatest),
						resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-2"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "3"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.nodes.#", "3"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.id"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),

						resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.#", "1"),
						resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.min", "1"),
						resource.TestCheckResourceAttr(resourceClusterName, "pool.0.autoscaler.0.max", "5"),
					),
				},
			},
		})
	})
}

func TestAccDataSourceLKECluster_controlPlane(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 2, func(tRetry *acceptance.TRetry) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(tRetry, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataControlPlane(t, clusterName, k8sVersionLatest, testRegion, true),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
						resource.TestCheckResourceAttr(dataSourceClusterName, "region", testRegion),
						resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", k8sVersionLatest),
						resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-2"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.nodes.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.autoscaler.#", "0"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.high_availability", "true"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.id"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),
					),
				},
			},
		})
	})
}
