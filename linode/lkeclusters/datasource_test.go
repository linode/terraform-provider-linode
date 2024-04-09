//go:build integration || lkeclusters

package lkeclusters_test

import (
	"context"
	"log"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeclusters/tmpl"
)

var (
	k8sVersionLatest string
	testRegion       string
)

func init() {
	// Get valid K8s versions for testing
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	versions, err := client.ListLKEVersions(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	k8sVersions := make([]string, len(versions))
	for i, v := range versions {
		k8sVersions[i] = v.ID
	}

	sort.Strings(k8sVersions)

	if len(k8sVersions) < 1 {
		log.Fatal("no k8s versions found")
	}

	k8sVersionLatest = k8sVersions[len(k8sVersions)-1]

	region, err := acceptance.GetRandomRegionWithCaps([]string{"kubernetes"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceLKEClusters_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_lke_clusters.test"

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
						acceptance.CheckResourceAttrGreaterThan(dataSourceName, "lke_clusters.#", 1),
					),
				},
				{
					Config: tmpl.DataFilter(t, clusterName, k8sVersionLatest, testRegion),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.label", clusterName),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.region", testRegion),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.k8s_version", k8sVersionLatest),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.status", "ready"),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.tags.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "lke_clusters.0.control_plane.high_availability", "false"),
					),
				},
			},
		})
	})
}
