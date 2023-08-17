//go:build unit

package helper_test

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestGetLatestVersion(t *testing.T) {
	versions := []string{"2.7.3", "0.5.0", "1.7.3", "8.0.27", "4.36.0", "2.56.100", "7.55.3", "8.0.26"}
	versionsMaps := make([]map[string]interface{}, len(versions))
	for i, ver := range versions {
		versionsMaps[i] = map[string]interface{}{
			"version": ver,
		}
	}

	var f helper.FilterConfig

	latest, err := f.GetLatestVersion(versionsMaps)
	if err != nil {
		t.Fatal(err)
	}

	if latest["version"] != "8.0.27" {
		t.Fatalf("expected 8.0.27, got %s", latest["version"])
	}
}
