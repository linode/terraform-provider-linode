//go:build integration || nfsfilesystem

package nfsfilesystem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystem/tmpl"
)

const dataSourceName = "data.linode_nfs_filesystem.foobar"

func TestAccDataSourceNFSFilesystem_basic_smoke(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	testRegion := acceptance.GetNFSTestRegion(t)

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
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("protocol_versions"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(string(linodego.NFSProtocolVersionV4)),
						})),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("mount_target_ips"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("tags"), knownvalue.SetSizeExact(0)),
				},
			},
		},
	})
}
