//go:build integration || nfsquota

package nfsquota_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota/tmpl"
)

func TestAccDataSourceNFSQuota_basic(t *testing.T) {
	t.Parallel()

	spaceID, filesystemID, mountPath := setupFilesystem(t)
	path := mountPath

	dataName := "data.linode_nfs_quota.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, spaceID, filesystemID, path, 2147483648, 2000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("space_id"), knownvalue.Int64Exact(int64(spaceID))),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("filesystem_id"), knownvalue.Int64Exact(int64(filesystemID))),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("path"), knownvalue.StringExact(path)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(2147483648)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(2000)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}
