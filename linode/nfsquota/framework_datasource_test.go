//go:build integration || nfsquota

package nfsquota_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota/tmpl"
)

func TestAccDataSourceNFSQuota_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	fs := acceptance.GetNFSTestFilesystem(t)

	dataName := "data.linode_nfs_quota.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, fs.SpaceID, fs.ID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						dataName, tfjsonpath.New("space_id"),
						resourceName, tfjsonpath.New("space_id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						dataName, tfjsonpath.New("filesystem_id"),
						resourceName, tfjsonpath.New("filesystem_id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						dataName, tfjsonpath.New("path"),
						resourceName, tfjsonpath.New("path"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						dataName, tfjsonpath.New("max_capacity_bytes"),
						resourceName, tfjsonpath.New("max_capacity_bytes"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						dataName, tfjsonpath.New("max_file_count"),
						resourceName, tfjsonpath.New("max_file_count"),
						compare.ValuesSame(),
					),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("is_root"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}
