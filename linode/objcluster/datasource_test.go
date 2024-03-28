//go:build (integration || objcluster) && !optional && !long_running

package objcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/objcluster/tmpl"
)

func TestAccDataSourceObjectCluster_basic(t *testing.T) {
	t.Parallel()

	objectStorageClusterID := "us-east-1"
	region := "us-east"
	resourceName := "data.linode_object_storage_cluster.foobar"
	staticSiteDomain := "website-us-east-1.linodeobjects.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, objectStorageClusterID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "id", objectStorageClusterID),
					resource.TestCheckResourceAttr(resourceName, "static_site_domain", staticSiteDomain),
				),
			},
		},
	})
}
