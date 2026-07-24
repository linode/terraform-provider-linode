//go:build integration || util

// NOTE: This test file needs to be tagged as integration because the
// package accesses the Linode API during init.

package acceptance

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetTestClient_noURLOverride(t *testing.T) {
	expectedURL := "api.linode.com"
	expectedVersion := "v4beta"

	t.Setenv("LINODE_URL", "")
	t.Setenv("LINODE_API_VERSION", "")

	client, err := GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	// baseURL and apiVersion are a private fields of
	// linodego.Client, so we need to access them using reflection
	rClient := reflect.ValueOf(*client)
	baseURL := rClient.FieldByName("baseURL").String()
	apiVersion := rClient.FieldByName("apiVersion").String()

	if baseURL != expectedURL {
		t.Fatalf("expected base url to be %s, got %s", expectedURL, baseURL)
	}

	if apiVersion != expectedVersion {
		t.Fatalf("expected api version to be %s, got %s", expectedVersion, apiVersion)
	}
}

func TestGetTestClient_URLOverride(t *testing.T) {
	expectedURL := "foo.linode.com"
	expectedVersion := "v4"

	t.Setenv("LINODE_URL", expectedURL)
	t.Setenv("LINODE_API_VERSION", expectedVersion)

	client, err := GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	// baseURL and apiVersion are a private fields of
	// linodego.Client, so we need to access them using reflection
	rClient := reflect.ValueOf(*client)
	baseURL := rClient.FieldByName("baseURL").String()
	apiVersion := rClient.FieldByName("apiVersion").String()

	if baseURL != expectedURL {
		t.Fatalf("expected base url to be %s, got %s", expectedURL, baseURL)
	}

	if apiVersion != expectedVersion {
		t.Fatalf("expected api version to be %s, got %s", expectedVersion, apiVersion)
	}
}

func TestAnyOfTestCheckFunc(t *testing.T) {
	err := errors.New("")
	checkFuncs1 := AnyOfTestCheckFunc(
		func(s *terraform.State) error { return err },
		func(s *terraform.State) error { return nil },
		func(s *terraform.State) error { return err },
	)
	checkFuncs2 := AnyOfTestCheckFunc(
		func(s *terraform.State) error { return nil },
	)

	checkFuncs3 := AnyOfTestCheckFunc(
		func(s *terraform.State) error { return err },
	)

	assert.NoError(t, checkFuncs1(nil))
	assert.NoError(t, checkFuncs2(nil))
	assert.Error(t, checkFuncs3(nil))
}

func TestRegionSelector_LocalFilters(t *testing.T) {
	selector := &RegionSelector{
		regions: []linodego.Region{
			{ID: "core-with-linodes", SiteType: "core", Capabilities: []string{"Linodes", "VPCs"}},
			{ID: "distributed-with-linodes", SiteType: "distributed", Capabilities: []string{"LINODES"}},
			{ID: "core-without-linodes", SiteType: "core", Capabilities: []string{"VPCs"}},
		},
	}

	filtered := selector.
		WithCapabilities(linodego.CapabilityLinodes).
		WithSiteType("CORE")

	filteredIDs, err := filtered.IDs()
	assert.NoError(t, err)
	assert.Equal(t, []string{"core-with-linodes"}, filteredIDs)

	originalIDs, err := selector.IDs()
	assert.NoError(t, err)
	assert.Len(t, originalIDs, 3)
	assert.Equal(t, "core-with-linodes", assertFirst(t, filtered))
}

func TestRegionSelector_Empty(t *testing.T) {
	selector := &RegionSelector{}

	_, err := selector.First()
	assert.EqualError(t, err, "no region matches the requested criteria")

	_, err = selector.Random()
	assert.EqualError(t, err, "no region matches the requested criteria")
}

func TestRegionSelector_Error(t *testing.T) {
	selector := &RegionSelector{err: errors.New("failed to select regions")}

	_, err := selector.
		WithCapabilities(linodego.CapabilityLinodes).
		FilteredBy(func(linodego.Region) bool { return true }).
		Random()
	assert.EqualError(t, err, "failed to select regions")
}

func TestRegionSelector_AvailableForAccountRejectsUnsupportedCapability(t *testing.T) {
	selector := &RegionSelector{}

	_, err := selector.
		AvailableForAccount(t.Context(), linodego.CapabilityCloudFirewall).
		Random()
	assert.EqualError(
		t,
		err,
		`capability "Cloud Firewall" is not supported by the account availability API`,
	)
}

func assertFirst(t *testing.T, selector *RegionSelector) string {
	t.Helper()

	region, err := selector.First()
	assert.NoError(t, err)
	return region
}
