//go:build integration || objglobalquotas

package objglobalquotas_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/objglobalquotas/tmpl"
)

func TestAccDataSourceObjGlobalQuotas_basic(t *testing.T) {
	t.Parallel()

	const dsAll = "data.linode_object_storage_global_quotas.all"
	const dsByQuotaID = "data.linode_object_storage_global_quotas.by-quota-id"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	quotas, err := client.ListObjectStorageGlobalQuotas(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list global quotas: %s", err)
	}

	if len(quotas) < 1 {
		t.Skipf("No available Object Storage global quota for testing. Skipping now...")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, quotas[0].QuotaID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas"),
						knownvalue.ListSizeExact(len(quotas)),
					),
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
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsAll,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("has_usage"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dsByQuotaID,
						tfjsonpath.New("quotas"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						dsByQuotaID,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_id"),
						knownvalue.StringExact(quotas[0].QuotaID),
					),
					statecheck.ExpectKnownValue(
						dsByQuotaID,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_type"),
						knownvalue.StringExact(quotas[0].QuotaType),
					),
					statecheck.ExpectKnownValue(
						dsByQuotaID,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("has_usage"),
						knownvalue.Bool(quotas[0].HasUsage),
					),
				},
			},
		},
	})
}
