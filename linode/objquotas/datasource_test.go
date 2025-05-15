//go:build integration || objquotas

package objquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
					// Filter and list object storage quotas match the endpoint type: E0
					acceptance.CheckResourceAttrGreaterThan(dsByEndpointType, "quotas.#", 2),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					// Check the first element of the Object Storage quotas
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("endpoint_type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("s3_endpoint"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_limit"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("resource_metric"),
						knownvalue.NotNull(),
					),

					// Filter and check the object storage quota match the endpoint type: E0
					statecheck.ExpectKnownValue(
						dsByEndpointType,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("endpoint_type"),
						knownvalue.StringExact("E0"),
					),
				},
			},
		},
	})
}
