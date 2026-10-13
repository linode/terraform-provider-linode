//go:build integration || nfsfilesystems

package nfsfilesystems_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystems/tmpl"
)

// minCapacityBytes is the minimum provisionable filesystem capacity (1 GiB).
const minCapacityBytes int64 = 1073741824

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatalf("failed to get random region: %s", err)
	}
	testRegion = region
}

func TestAccDataSourceNFSFilesystems_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nfs_filesystems.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_filesystems").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_filesystems").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.StringExact(testRegion)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_filesystems").AtSliceIndex(0).AtMapKey("max_capacity_bytes"),
						knownvalue.Int64Exact(minCapacityBytes)),
				},
			},
		},
	})
}

func TestAccDataSourceNFSFilesystems_filterLabel(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nfs_filesystems.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterLabel(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_filesystems").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label)),
				},
			},
		},
	})
}
