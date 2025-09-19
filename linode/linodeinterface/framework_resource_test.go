//go:build integration || linodeinterface

package linodeinterface_test

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
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	linodeinstancetmpl "github.com/linode/terraform-provider-linode/v3/linode/acceptance/tmpl"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/linodeinterface/tmpl"
)

const testInterfaceResName = "linode_interface.test"

var testRegion string

func init() {
	resource.AddTestSweepers("linode_interface", &resource.Sweeper{
		Name: "linode_interface",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes, linodego.CapabilityVlans, linodego.CapabilityVPCs}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err)
	}

	// Get all instances and sweep their interfaces
	instances, err := client.ListInstances(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get instances: %s", err)
	}

	for _, instance := range instances {
		if !acceptance.ShouldSweep(prefix, instance.Label) {
			continue
		}

		// Get instance configs to find interfaces
		configs, err := client.ListInstanceConfigs(context.Background(), instance.ID, nil)
		if err != nil {
			continue // Skip if we can't get configs
		}

		// Delete non-primary interfaces from configs
		for _, config := range configs {
			for i, iface := range config.Interfaces {
				// Skip eth0 (primary interface) and other essential interfaces
				if i == 0 || iface.Purpose == linodego.InterfacePurposePublic {
					continue
				}

				// For sweep purposes, we'll let the instance deletion handle interface cleanup
				// since interfaces are tied to instances
			}
		}
	}

	return nil
}

func TestAccLinodeInterface_vlan_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	vlanLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VLANBasic(t, label, testRegion, vlanLabel, "192.168.200.5/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vlan").AtMapKey("vlan_label"), knownvalue.StringExact(vlanLabel)),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("vlan").AtMapKey("ipam_address"),
						knownvalue.StringExact("192.168.200.5/24"),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicBasic(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("primary"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("assigned_addresses"),
						knownvalue.NotNull(),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_ipv4_ipv6(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicWithIPv4AndIPv6(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("primary"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.StringExact("/64"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("assigned_addresses"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("assigned_ranges"),
						knownvalue.NotNull(),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_ipv6_only(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicWithIPv6(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					// Verify IPv4 addresses is explicitly empty (not omitted)
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(0),
					),
					// Verify IPv6 ranges are configured
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.StringExact("/64"),
					),
					// Verify no IPv4 addresses are assigned due to empty list
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("assigned_addresses"),
						knownvalue.ListSizeExact(0),
					),
					// Verify IPv6 ranges are assigned
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("assigned_ranges"),
						knownvalue.NotNull(),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_update_addresses(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicWithIPv4(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("primary"),
						knownvalue.Bool(true),
					),
				},
				Check: checkInterfaceExists,
			},
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicUpdatedIPv4(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(2),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("primary"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(1).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(1).AtMapKey("primary"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges"),
						knownvalue.ListSizeExact(2),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges").AtSliceIndex(0).AtMapKey("range"),
						knownvalue.StringExact("/64"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("public").AtMapKey("ipv6").AtMapKey("ranges").AtSliceIndex(1).AtMapKey("range"),
						knownvalue.StringExact("/64"),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_vpc_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCBasic(t, label, testRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_vpc_with_ipv4(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCWithIPv4(t, label, testRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("addresses"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact("auto"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("addresses").AtSliceIndex(0).AtMapKey("primary"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName,
						tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("assigned_addresses"),
						knownvalue.NotNull(),
					),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_vpc_update_ipv4(t *testing.T) {
	t.Skip("Backend issue prevents adding additional IPv4 addresses to VPC interfaces")
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCWithIPv4(t, label, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					checkInterfaceExists,
					resource.TestCheckResourceAttrSet(testInterfaceResName, "id"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "linode_id"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "vpc.subnet_id"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.#", "1"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.0.address", "auto"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.0.primary", "true"),
				),
			},
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCUpdatedIPv4(t, label, testRegion, "10.0.0.0/24"),
				Check: resource.ComposeTestCheckFunc(
					checkInterfaceExists,
					resource.TestCheckResourceAttrSet(testInterfaceResName, "id"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "linode_id"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "vpc.subnet_id"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.#", "2"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.0.address", "auto"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.0.primary", "true"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.1.address", "auto"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.addresses.1.primary", "false"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.ranges.#", "1"),
					resource.TestCheckResourceAttr(testInterfaceResName, "vpc.ipv4.ranges.0.range", "/32"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "vpc.ipv4.assigned_addresses.#"),
					resource.TestCheckResourceAttrSet(testInterfaceResName, "vpc.ipv4.assigned_ranges.#"),
				),
			},
		},
	})
}

func TestAccLinodeInterface_public_default_route(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicDefaultRouteIPv6(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("default_route").AtMapKey("ipv4"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("default_route").AtMapKey("ipv6"), knownvalue.Bool(true)),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_vpc_default_route(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCDefaultRouteIPv4(t, label, testRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("default_route").AtMapKey("ipv4"), knownvalue.Bool(true)),
					// VPC interfaces don't support IPv6, so we don't expect ipv6 field to be set
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

// Smoke test to run a subset of tests
func TestSmokeTests_interface(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"TestAccLinodeInterface_vlan_basic", TestAccLinodeInterface_vlan_basic},
		{"TestAccLinodeInterface_public_basic", TestAccLinodeInterface_public_basic},
		{"TestAccLinodeInterface_public_ipv6_only", TestAccLinodeInterface_public_ipv6_only},
		{"TestAccLinodeInterface_vpc_basic", TestAccLinodeInterface_vpc_basic},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// Helper function to check if interface exists
func checkInterfaceExists(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_interface" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing Interface ID %v to int", rs.Primary.ID)
		}

		linodeID, err := strconv.Atoi(rs.Primary.Attributes["linode_id"])
		if err != nil {
			return fmt.Errorf("Error parsing Linode ID %v to int", rs.Primary.Attributes["linode_id"])
		}

		// Use second generation interface API to get the interface directly
		_, err = client.GetInterface(context.Background(), linodeID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving interface %d for instance %d: %s", id, linodeID, err)
		}
	}

	return nil
}

// Helper function to check if interface is destroyed
func checkInterfaceDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_interface" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing Interface ID %v to int", rs.Primary.ID)
		}

		linodeID, err := strconv.Atoi(rs.Primary.Attributes["linode_id"])
		if err != nil {
			return fmt.Errorf("Error parsing Linode ID %v to int", rs.Primary.Attributes["linode_id"])
		}

		if id == 0 {
			// Don't try to delete interface 0 (primary interface)
			continue
		}

		// Use second generation interface API to check if interface still exists
		_, err = client.GetInterface(context.Background(), linodeID, id)
		if err != nil {
			if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
				// Interface doesn't exist, which is expected after destroy
				continue
			}
			return fmt.Errorf("Error checking interface %d for instance %d: %s", id, linodeID, err)
		}

		// If we get here, the interface still exists
		return fmt.Errorf("Interface with id %d still exists", id)
	}

	return nil
}

func TestAccLinodeInterface_vpc_default_ip(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCDefaultIP(t, label, testRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_default_ip(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicDefaultIP(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_public_empty_ip_objects(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.PublicEmptyIPObjects(t, label, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					// Note: addresses may be nil when empty, so just check the fields exist
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("public").AtMapKey("ipv4"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("public").AtMapKey("ipv6"), knownvalue.NotNull()),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}

func TestAccLinodeInterface_vpc_empty_ip_objects(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCEmptyIPObjects(t, label, testRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
					// Note: addresses may be nil when empty, so just check the ipv4 field exists
					// VPC interfaces don't support IPv6, so we don't check for ipv6
					statecheck.ExpectKnownValue(testInterfaceResName, tfjsonpath.New("vpc").AtMapKey("ipv4"), knownvalue.NotNull()),
				},
				Check: checkInterfaceExists,
			},
		},
	})
}
