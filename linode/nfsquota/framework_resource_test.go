//go:build integration || nfsquota

package nfsquota_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota/tmpl"
)

const resourceName = "linode_nfs_quota.foobar"

// filesystemMountPath extracts the exported mount path root from a filesystem.
// The backend derives this path from the filesystem label plus a generated
// suffix (e.g. "/my-label-1c8"), and NFS Quota paths must be prefixed with it.
func filesystemMountPath(t *testing.T, fs *linodego.NFSFilesystem) string {
	t.Helper()

	if fs.MountTargetFQDN != nil {
		if _, path, found := strings.Cut(*fs.MountTargetFQDN, ":"); found {
			return path
		}
	}

	for _, target := range fs.MountTargetIPs {
		if idx := strings.LastIndex(target, ":/"); idx != -1 {
			return target[idx+1:]
		}
	}

	t.Fatalf("unable to determine mount path for filesystem %d", fs.ID)
	return ""
}

// setupFilesystem provisions an NFS Space and Filesystem for use as a parent
// of the NFS Quota under test. It registers cleanup via t.Cleanup and returns
// the resulting space ID, filesystem ID, and the exported mount path root that
// NFS Quota paths must be prefixed with.
func setupFilesystem(t *testing.T) (spaceID, filesystemID int, mountPath string) {
	t.Helper()

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	//region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	//if err != nil {
	//	t.Fatalf("failed to get region: %s", err)
	//}
	region := "pl-labkrk-2"

	ctx := context.Background()
	label := acctest.RandomWithPrefix("tf-test")

	space, err := client.CreateNFSSpace(ctx, linodego.NFSSpaceCreateOptions{
		Label: label,
	})
	if err != nil {
		t.Fatalf("failed to create NFS space: %s", err)
	}
	t.Cleanup(func() {
		if err := client.DeleteNFSSpace(context.Background(), space.ID); err != nil {
			if !linodego.IsNotFound(err) {
				t.Errorf("failed to delete NFS space %d: %s", space.ID, err)
			}
		}
	})

	if _, err := client.WaitForNFSSpaceStatus(ctx, space.ID, linodego.NFSSpaceStatusActive); err != nil {
		t.Fatalf("failed to wait for NFS space: %s", err)
	}

	fs, err := client.CreateNFSFilesystem(ctx, space.ID, linodego.NFSFilesystemCreateOptions{
		Label:            label,
		Region:           region,
		MaxCapacityBytes: 10 * 1024 * 1024 * 1024,
	})
	if err != nil {
		t.Fatalf("failed to create NFS filesystem: %s", err)
	}
	t.Cleanup(func() {
		if err := client.DeleteNFSFilesystem(context.Background(), space.ID, fs.ID); err != nil {
			if !linodego.IsNotFound(err) {
				t.Errorf("failed to delete NFS filesystem %d: %s", fs.ID, err)
			}
		}
	})

	fs, err = client.WaitForNFSFilesystemStatus(ctx, space.ID, fs.ID, linodego.NFSFilesystemStatusActive)
	if err != nil {
		t.Fatalf("failed to wait for NFS filesystem: %s", err)
	}

	return space.ID, fs.ID, filesystemMountPath(t, fs)
}

func TestAccResourceNFSQuota_basic(t *testing.T) {
	t.Parallel()

	spaceID, filesystemID, mountPath := setupFilesystem(t)
	path := mountPath

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, spaceID, filesystemID, path, 2147483648, 2000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("space_id"), knownvalue.Int64Exact(int64(spaceID))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("filesystem_id"), knownvalue.Int64Exact(int64(filesystemID))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("path"), knownvalue.StringExact(path)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(2147483648)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(2000)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return fmt.Sprintf("%d,%d,%s", spaceID, filesystemID, rs.Primary.ID), nil
				},
			},
		},
	})
}

func TestAccResourceNFSQuota_update(t *testing.T) {
	t.Parallel()

	spaceID, filesystemID, mountPath := setupFilesystem(t)
	path := mountPath

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkNFSQuotaDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, spaceID, filesystemID, path, 2147483648, 2000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(2147483648)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(2000)),
				},
			},
			{
				Config: tmpl.Updates(t, spaceID, filesystemID, path, 4294967296, 4000),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_capacity_bytes"), knownvalue.Int64Exact(4294967296)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("max_file_count"), knownvalue.Int64Exact(4000)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("default_user_limit").AtMapKey("max_capacity_bytes"),
						knownvalue.Int64Exact(1073741824),
					),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("user_group_config").AtMapKey("user_limits").AtSliceIndex(0).AtMapKey("identifier"),
						knownvalue.StringExact("1001"),
					),
				},
			},
		},
	})
}

func checkNFSQuotaDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nfs_quota" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing NFS Quota ID %v: %s", rs.Primary.ID, err)
		}

		spaceID, err := strconv.Atoi(rs.Primary.Attributes["space_id"])
		if err != nil {
			return fmt.Errorf("error parsing space_id %v: %s", rs.Primary.Attributes["space_id"], err)
		}

		filesystemID, err := strconv.Atoi(rs.Primary.Attributes["filesystem_id"])
		if err != nil {
			return fmt.Errorf("error parsing filesystem_id %v: %s", rs.Primary.Attributes["filesystem_id"], err)
		}

		_, err = client.GetNFSQuota(context.Background(), spaceID, filesystemID, id)
		if err == nil {
			return fmt.Errorf("NFS Quota with ID %d still exists", id)
		}
		if !linodego.IsNotFound(err) {
			return fmt.Errorf("error getting NFS Quota with ID %d: %s", id, err)
		}
	}

	return nil
}
