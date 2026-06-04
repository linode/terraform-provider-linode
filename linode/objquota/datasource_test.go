//go:build integration || objquota

package objquota_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/objquota/tmpl"
)

func TestAccDataSourceObjQuota_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_object_storage_quota.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	quotas, err := client.ListObjectStorageQuotas(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list quotas: %s", err)
	}

	if len(quotas) < 1 {
		t.Skipf("No available Object Storage quota for testing. Skipping now...")
	}

	selectedQuota := quotas[0]
	for _, quota := range quotas {
		if quota.HasUsage {
			selectedQuota = quota
			break
		}
	}

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("id"),
			knownvalue.StringExact(selectedQuota.QuotaID),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("quota_id"),
			knownvalue.StringExact(selectedQuota.QuotaID),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("quota_name"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("endpoint_type"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("s3_endpoint"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("description"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("quota_limit"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("resource_metric"),
			knownvalue.NotNull(),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("quota_type"),
			knownvalue.StringExact(selectedQuota.QuotaType),
		),
		statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("has_usage"),
			knownvalue.Bool(selectedQuota.HasUsage),
		),
	}

	if selectedQuota.HasUsage {
		checks = append(checks,
			statecheck.ExpectKnownValue(
				resourceName,
				tfjsonpath.New("quota_usage").AtMapKey("quota_limit"),
				knownvalue.NotNull(),
			),
			statecheck.ExpectKnownValue(
				resourceName,
				tfjsonpath.New("quota_usage").AtMapKey("usage"),
				knownvalue.NotNull(),
			),
		)
	} else {
		checks = append(checks,
			statecheck.ExpectKnownValue(
				resourceName,
				tfjsonpath.New("quota_usage"),
				knownvalue.Null(),
			),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            tmpl.DataBasic(t, selectedQuota.QuotaID),
				ConfigStateChecks: checks,
			},
		},
	})
}
