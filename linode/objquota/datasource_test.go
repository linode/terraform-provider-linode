//go:build integration || objquota

package objquota_test

import (
	"context"
	"fmt"
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
		fmt.Errorf("Error getting client: %s", err.Error())
	}

	quotas, err := client.ListObjectStorageQuotas(context.Background(), nil)
	if err != nil {
		fmt.Errorf("Error listing quotas: %s", err.Error())
	}

	if len(quotas) < 1 {
		t.Skipf("No available Object Storage quota for testing. Skipping now...")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, quotas[0].QuotaID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("id"),
						knownvalue.StringExact(quotas[0].QuotaID),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("quota_id"),
						knownvalue.StringExact(quotas[0].QuotaID),
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
						tfjsonpath.New("quota_usage").AtMapKey("quota_limit"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("quota_usage").AtMapKey("usage"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
