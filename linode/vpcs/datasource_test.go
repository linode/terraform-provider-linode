//go:build integration || vpcs

package vpcs_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcs/tmpl"
)

// preConfigVPCPoll returns a PreConfig function that waits for a VPC with the
// given label to be returned by the list endpoint.
func preConfigVPCPoll(t testing.TB, vpcLabel string) func() {
	return func() {
		client, err := acceptance.GetTestClient()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := waitForVPCWithLabel(client, vpcLabel, 60); err != nil {
			t.Fatal(err)
		}
	}
}

// waitForVPCWithLabel polls the list endpoint until a VPC with the given label
// is returned.
func waitForVPCWithLabel(client *linodego.Client, label string, timeoutSeconds int) (*linodego.VPC, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSeconds))
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			vpcs, err := client.ListVPCs(ctx,
				&linodego.ListOptions{Filter: fmt.Sprintf("{\"label\": \"%s\"}", label)})
			if err != nil {
				return nil, err
			}

			for i := range vpcs {
				if vpcs[i].Label == label {
					return &vpcs[i], nil
				}
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for VPC %s: %w", label, ctx.Err())
		}
	}
}

func TestAccDataSourceVPCs_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Error(fmt.Errorf("failed to get region with VPC capability: %w", err))
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpcs.#", 0),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.description"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.updated"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCs_dualStack(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityVPCs,
		linodego.CapabilityVPCDualStack,
	}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataDualStack(t, vpcLabel, targetRegion),
			},
			{
				PreConfig: preConfigVPCPoll(t, vpcLabel),
				Config:    tmpl.DataDualStack(t, vpcLabel, targetRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("created"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("updated"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv6"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv6").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceVPCs_filterByLabel(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterLabel(t, vpcLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "vpcs.#", 0),
					acceptance.CheckResourceAttrContains(resourceName, "vpcs.0.label", "tf-test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcs.0.updated"),
				),
			},
		},
	})
}

func TestAccDataSourceVPCs_ipv4(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_vpcs.foobar"
	vpcLabel := acctest.RandomWithPrefix("tf-test")

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityVPCs,
		linodego.CapabilityVPCCustomIPv4Ranges,
	}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataIPv4(t, vpcLabel, targetRegion, "10.0.0.0/8"),
			},
			{
				PreConfig: preConfigVPCPoll(t, vpcLabel),
				Config:    tmpl.DataIPv4(t, vpcLabel, targetRegion, "10.0.0.0/8"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("created"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("updated"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("vpcs").AtSliceIndex(0).AtMapKey("ipv4").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.StringExact("10.0.0.0/8"),
					),
				},
			},
		},
	})
}
