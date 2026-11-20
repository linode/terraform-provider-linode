package acceptance

import (
	"os"
	"testing"

	"github.com/linode/linodego/v2"
)

const nfsTestRegionEnvVar = "ACC_NFS_TEST_REGION"

// GetNFSTestRegion returns the region to run NFS tests in. The region can be
// overridden with the ACC_NFS_TEST_REGION environment variable; otherwise a
// region with the NFS Storage capability is selected. The test is skipped if
// no such region is available.
func GetNFSTestRegion(t testing.TB) string {
	t.Helper()

	if region := os.Getenv(nfsTestRegionEnvVar); region != "" {
		return region
	}

	region, err := GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityNFSStorage},
		"any",
	)
	if err != nil {
		t.Skipf(
			"skipping NFS test; no region with the %q capability found (%s). Set %s to override",
			linodego.CapabilityNFSStorage, err, nfsTestRegionEnvVar,
		)
	}

	return region
}
