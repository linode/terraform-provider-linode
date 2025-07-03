//go:build integration || lke

package lke_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/lke/tmpl"
)

const dataSourceClusterName = "data.linode_lke_cluster.test"

func TestAccDataSourceLKECluster_taints_labels(t *testing.T) {
	t.Parallel()

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataTaintsLabels(
						t, clusterName, k8sVersionLatest, testRegion, []tmpl.TaintData{
							{
								Effect: "PreferNoSchedule",
								Key:    "foo",
								Value:  "bar",
							},
						},
						map[string]string{"foo": "bar"},
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.nodes.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.taints.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.labels.foo", "bar"),
					),
				},
			},
		})
	})
}

func TestAccDataSourceLKECluster_basic(t *testing.T) {
	t.Parallel()

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, clusterName, k8sVersionLatest, testRegion),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
						resource.TestCheckResourceAttr(dataSourceClusterName, "region", testRegion),
						resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", k8sVersionLatest),
						resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tier", "standard"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "3"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.disk_encryption"),
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

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-1"),
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

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		testIPv4 := "0.0.0.0/0"
		testIPv6 := "2001:db8::/32"

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataControlPlane(t, clusterName, k8sVersionLatest, testRegion, testIPv4, testIPv6, false, true),
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
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.high_availability", "false"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.acl.0.enabled", "true"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.acl.0.addresses.0.ipv4.0", testIPv4),
						resource.TestCheckResourceAttr(dataSourceClusterName, "control_plane.0.acl.0.addresses.0.ipv6.0", testIPv6),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "pools.0.id"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),
					),
				},
			},
		})
	})
}

func TestAccDataSourceLKECluster_enterprise(t *testing.T) {
	t.Parallel()

	if k8sVersionEnterprise == "" {
		t.Skip("No available k8s version for LKE Enterprise test. Skipping now...")
	}

	enterpriseRegion, err := acceptance.GetRandomRegionWithCaps([]string{"Kubernetes Enterprise"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	client, err := acceptance.GetTestClient()

	enterpriseVersions, err := client.ListLKETierVersions(context.Background(), "enterprise", nil)
	if err != nil {
		log.Fatal(err)
	}

	var k8sVersionEnterprise string
	if len(enterpriseVersions) < 1 {
		t.Skip("No available k8s version for LKE Enterprise test. Skipping now...")
	} else {
		k8sVersionEnterprise = enterpriseVersions[0].ID
	}

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		clusterName := acctest.RandomWithPrefix("tf_test")
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataEnterprise(t, clusterName, k8sVersionEnterprise, enterpriseRegion, "on_recycle"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceClusterName, "label", clusterName),
						resource.TestCheckResourceAttr(dataSourceClusterName, "region", enterpriseRegion),
						resource.TestCheckResourceAttr(dataSourceClusterName, "k8s_version", k8sVersionEnterprise),
						resource.TestCheckResourceAttr(dataSourceClusterName, "status", "ready"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tier", "enterprise"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.#", "1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.type", "g6-standard-1"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.count", "3"),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.k8s_version", k8sVersionEnterprise),
						resource.TestCheckResourceAttr(dataSourceClusterName, "pools.0.update_strategy", "on_recycle"),
						resource.TestCheckResourceAttrSet(dataSourceClusterName, "kubeconfig"),
					),
				},
			},
		})
	})
}
