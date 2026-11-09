//go:build integration || nfsfilesystem

package nfsfilesystem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystem/tmpl"
)

const dataSourceName = "data.linode_nfs_filesystem.foobar"

func TestAccDataSourceNFSFilesystem_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("region"), knownvalue.StringExact(testRegion)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(minCapacityBytes)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("space_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
