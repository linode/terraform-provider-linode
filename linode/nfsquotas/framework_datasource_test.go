//go:build integration || nfsquotas

package nfsquotas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquotas/tmpl"
)

func TestAccDataSourceNFSQuotas_basic(t *testing.T) {
	t.Parallel()

	fs := acceptance.GetNFSTestFilesystem(t)
	dataSourceName := "data.linode_nfs_quotas.foobar"

	quota := tfjsonpath.New("nfs_quotas").AtSliceIndex(0)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, fs.SpaceID, fs.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("nfs_quotas"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("space_id"), knownvalue.Int64Exact(int64(fs.SpaceID))),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("filesystem_id"), knownvalue.Int64Exact(int64(fs.ID))),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("is_root"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("path"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("max_capacity_bytes"), knownvalue.Int64Exact(fs.MaxCapacityBytes)),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("max_file_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, quota.AtMapKey("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}
