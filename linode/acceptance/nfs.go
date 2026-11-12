package acceptance

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/linode/linodego/v2"
)

const (
	nfsTestRegionEnvVar = "ACC_NFS_TEST_REGION"

	// NFSTestFilesystemMinCapacityBytes is the minimum capacity for a reused
	// test filesystem.
	NFSTestFilesystemMinCapacityBytes int64 = 10737418240
)

var (
	nfsTestFilesystemOnce sync.Once
	nfsTestFilesystem     *linodego.NFSFilesystem
	nfsTestFilesystemErr  error
)

// GetNFSTestRegion returns the region to run NFS tests in, skipping the test
// if the region is not provided.
func GetNFSTestRegion(t testing.TB) string {
	t.Helper()

	region := os.Getenv(nfsTestRegionEnvVar)
	if region == "" {
		t.Skipf("skipping NFS test; set the %s environment variable to run", nfsTestRegionEnvVar)
	}

	return region
}

// GetNFSTestFilesystem returns an active filesystem in the NFS test region,
// creating a space and filesystem if the account has none. Deleting a subpath
// quota leaves its directory behind, which blocks filesystem deletion through
// the API, so the space and filesystem are shared across tests and must be
// cleaned up manually.
func GetNFSTestFilesystem(t testing.TB) *linodego.NFSFilesystem {
	t.Helper()

	region := GetNFSTestRegion(t)

	nfsTestFilesystemOnce.Do(func() {
		nfsTestFilesystem, nfsTestFilesystemErr = findOrCreateNFSTestFilesystem(region)
	})

	if nfsTestFilesystemErr != nil {
		t.Fatalf("failed to get NFS filesystem for testing: %s", nfsTestFilesystemErr)
	}

	return nfsTestFilesystem
}

func findOrCreateNFSTestFilesystem(region string) (*linodego.NFSFilesystem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client, err := GetTestClient()
	if err != nil {
		return nil, err
	}

	spaces, err := client.ListNFSSpaces(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list NFS spaces: %w", err)
	}

	for _, space := range spaces {
		if space.Status != linodego.NFSSpaceStatusActive {
			continue
		}

		filesystems, err := client.ListNFSFilesystems(ctx, space.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list NFS filesystems in space %d: %w", space.ID, err)
		}

		for _, fs := range filesystems {
			if fs.Region == region &&
				fs.Status == linodego.NFSFilesystemStatusActive &&
				fs.MaxCapacityBytes >= NFSTestFilesystemMinCapacityBytes {
				return &fs, nil
			}
		}
	}

	label := acctest.RandomWithPrefix("tf-test-nfs")

	space, err := client.CreateNFSSpace(ctx, linodego.NFSSpaceCreateOptions{
		Label: label + "-space",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create NFS space: %w", err)
	}

	if _, err := client.WaitForNFSSpaceStatus(ctx, space.ID, linodego.NFSSpaceStatusActive); err != nil {
		return nil, err
	}

	fs, err := client.CreateNFSFilesystem(ctx, space.ID, linodego.NFSFilesystemCreateOptions{
		Label:            label,
		Region:           region,
		MaxCapacityBytes: NFSTestFilesystemMinCapacityBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create NFS filesystem in space %d: %w", space.ID, err)
	}

	return client.WaitForNFSFilesystemStatus(ctx, space.ID, fs.ID, linodego.NFSFilesystemStatusActive)
}
