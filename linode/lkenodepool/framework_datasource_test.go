//go:build integration || lkenodepool

package lkenodepool_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/lkenodepool/tmpl"
)

func TestAccDataSourceLKENodePool_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_lke_node_pool.test"
	clusterLabel := acctest.RandomWithPrefix("tf_test_")
	poolTag := acctest.RandomWithPrefix("tf_test_")

	templateData := createTemplateData()
	templateData.ClusterLabel = clusterLabel
	templateData.PoolTag = poolTag
	templateData.AutoscalerEnabled = true
	templateData.AutoscalerMin = 1
	templateData.AutoscalerMax = 2
	templateData.NodeCount = 1
	templateData.Labels = map[string]string{"foo": "bar"}
	templateData.Taints = []tmpl.TaintData{
		{
			Effect: "PreferNoSchedule",
			Key:    "foo",
			Value:  "bar",
		},
	}

	acceptance.RunTestWithRetries(t, 2, func(t *acceptance.WrappedT) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			CheckDestroy:             acceptance.CheckLKEClusterDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, &templateData),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						resource.TestCheckResourceAttrSet(dataSourceName, "cluster_id"),
						resource.TestCheckResourceAttr(dataSourceName, "disk_encryption", "enabled"),
						resource.TestCheckResourceAttr(dataSourceName, "type", "g6-standard-1"),
						resource.TestCheckResourceAttr(dataSourceName, "node_count", "1"),
						resource.TestCheckNoResourceAttr(dataSourceName, "firewall_id"),
						resource.TestCheckNoResourceAttr(dataSourceName, "k8s_version"),
						resource.TestCheckNoResourceAttr(dataSourceName, "update_strategy"),

						resource.TestCheckResourceAttr(dataSourceName, "autoscaler.enabled", "true"),
						resource.TestCheckResourceAttr(dataSourceName, "autoscaler.min", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "autoscaler.max", "2"),

						resource.TestCheckResourceAttr(dataSourceName, "nodes.#", "1"),
						resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.id"),
						resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.instance_id"),
						resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.status"),

						resource.TestCheckResourceAttr(dataSourceName, "disks.#", "0"),

						resource.TestCheckResourceAttr(dataSourceName, "labels.%", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "labels.foo", "bar"),

						resource.TestCheckResourceAttr(dataSourceName, "taints.#", "1"),
						resource.TestCheckResourceAttr(dataSourceName, "taints.0.effect", "PreferNoSchedule"),
						resource.TestCheckResourceAttr(dataSourceName, "taints.0.key", "foo"),
						resource.TestCheckResourceAttr(dataSourceName, "taints.0.value", "bar"),

						resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
						resource.TestCheckResourceAttr(dataSourceName, "tags.0", "external"),
						resource.TestCheckResourceAttr(dataSourceName, "tags.1", poolTag),
					),
				},
				{
					Config:      tmpl.DataClusterNotFound(t, &templateData),
					ExpectError: regexp.MustCompile(`\[404\] Not found`),
				},
				{
					Config:      tmpl.DataNodePoolNotFound(t, &templateData),
					ExpectError: regexp.MustCompile(`\[404\] Not found`),
				},
			},
		})
	})
}
