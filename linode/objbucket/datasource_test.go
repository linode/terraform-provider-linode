//go:build integration || objbucket

package objbucket_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket/tmpl"
)

func TestAccDataSourceBucket_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_object_storage_bucket.foobar"
	objectStorageBucketName := acctest.RandomWithPrefix("tf-test")

	acceptance.RunTestRetry(t, 5, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkBucketDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, objectStorageBucketName, testCluster),
					Check: resource.ComposeTestCheckFunc(
						checkBucketExists,
						resource.TestCheckResourceAttr(resourceName, "cluster", testCluster),
						resource.TestCheckResourceAttr(resourceName, "label", objectStorageBucketName),
						resource.TestCheckResourceAttrSet(resourceName, "hostname"),
						resource.TestCheckResourceAttrSet(resourceName, "created"),
						resource.TestCheckResourceAttrSet(resourceName, "objects"),
						resource.TestCheckResourceAttrSet(resourceName, "size"),
					),
				},
			},
		})
	})
}
