package acceptance

import (
	"os"
	"testing"
)

const nfsTestRegionEnvVar = "ACC_NFS_TEST_REGION"

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
