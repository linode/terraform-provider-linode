//go:build integration || nfsfilesystem

package nfsfilesystem_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystem/tmpl"
)

const resourceName = "linode_nfs_filesystem.foobar"

// minCapacityBytes is the minimum provisionable filesystem capacity (1 GiB).
const minCapacityBytes int64 = 1073741824

var testRegion string

func init() {
	resource.AddTestSweepers("linode_nfs_filesystem", &resource.Sweeper{
		Name: "linode_nfs_filesystem",
		F:    sweep,
	})

	//region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityLinodes}, "core")
	//if err != nil {
	//	log.Fatalf("failed to get random region: %s", err)
	//}
	//testRegion = region
	// only use pl-labkrk-2 for now since other regions are not available for NFS Filesystem creation
	testRegion = "pl-labkrk-2"
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	spaces, err := client.ListNFSSpaces(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error getting NFS Spaces: %s", err)
	}

	for _, space := range spaces {
		filesystems, err := client.ListNFSFilesystems(context.Background(), space.ID, nil)
		if err != nil {
			return fmt.Errorf("error getting NFS Filesystems: %s", err)
		}

		for _, fs := range filesystems {
			if !acceptance.ShouldSweep(prefix, fs.Label) {
				continue
			}
			if err := client.DeleteNFSFilesystem(context.Background(), space.ID, fs.ID); err != nil {
				return fmt.Errorf("error destroying %v during sweep: %s", fs.Label, err)
			}
		}
	}

	return nil
}

func TestAccResourceNFSFilesystem_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSFilesystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("region"), knownvalue.StringExact(testRegion)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("status"), knownvalue.StringExact(string(linodego.NFSFilesystemStatusActive))),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(minCapacityBytes)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("space_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("mount_target_fqdn"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNFSFilesystem_update(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	newLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSFilesystemDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, minCapacityBytes),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
				},
			},
			{
				Config: tmpl.Updates(t, newLabel, testRegion, minCapacityBytes*2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(newLabel)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(minCapacityBytes*2)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("tags"), knownvalue.SetSizeExact(2)),
				},
			},
		},
	})
}

func checkNFSFilesystemDestroy(s *terraform.State) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nfs_filesystem" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetNFSFilesystem(context.Background(), id)
		if err == nil {
			return fmt.Errorf("NFS Filesystem with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting NFS Filesystem with id %d", id)
		}
	}

	return nil
}
