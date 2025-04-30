//go:build integration || objquota

package objquota_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/objquota/tmpl"
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, quotas[0].QuotaID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", quotas[0].QuotaID),
					resource.TestCheckResourceAttrSet(resourceName, "quota_name"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_type"),
					resource.TestCheckResourceAttrSet(resourceName, "s3_endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_metric"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_usage.quota_limit"),
					resource.TestCheckResourceAttrSet(resourceName, "quota_usage.usage"),
				),
			},
		},
	})

}
