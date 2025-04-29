//go:build integration || objquotas

package objquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/objquotas/tmpl"
)

func TestAccDataSourceObjQuotas_basic(t *testing.T) {
	t.Parallel()

	const dsAll = "data.linode_object_storage_quotas.all"
	const dsByQuotaName = "data.linode_object_storage_quotas.by-quota-name"
	const dsByEndpointType = "data.linode_object_storage_quotas.by-endpoint-type"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					// List all object storage quotas
					acceptance.CheckResourceAttrGreaterThan(dsAll, "quotas.#", 2),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.id"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.quota_name"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.endpoint_type"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.s3_endpoint"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.description"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dsAll, "quotas.0.resource_metric"),

					// Filter and list object storage quotas match the quota name: max_buckets
					acceptance.CheckResourceAttrGreaterThan(dsByQuotaName, "quotas.#", 2),
					resource.TestCheckResourceAttr(dsByQuotaName, "quotas.0.quota_name", "max_buckets"),

					// Filter and list object storage quotas match the endpoint type: E0
					acceptance.CheckResourceAttrGreaterThan(dsByEndpointType, "quotas.#", 2),
					resource.TestCheckResourceAttr(dsByEndpointType, "quotas.0.endpoint_type", "E0"),
				),
			},
		},
	})
}
