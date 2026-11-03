//go:build integration || natgatewayinterfaces

package natgatewayinterfaces_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	linodeinstancetmpl "github.com/linode/terraform-provider-linode/v4/linode/acceptance/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/natgatewayinterfaces/tmpl"
)

const (
	natGatewayResName      = "linode_natgateway.test"
	interfacesByGatewayRes = "data.linode_natgateway_interfaces.by_gateway"
	interfacesByAddressRes = "data.linode_natgateway_interfaces.by_address"
)

// waitForNATGatewayInterface polls the API until the NAT gateway lists at least
// one interface. The interface is attached asynchronously, so the data sources
// must not be read before it shows up.
func waitForNATGatewayInterface(resName string, timeout time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resName)
		}

		gatewayID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing NAT gateway ID %q: %w", rs.Primary.ID, err)
		}

		client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

		deadline := time.Now().Add(timeout)
		for time.Now().Before(deadline) {
			interfaces, err := client.NATGatewayListInterfaces(context.Background(), gatewayID, nil)
			if err != nil {
				return fmt.Errorf("error listing interfaces for NAT gateway %d: %w", gatewayID, err)
			}
			if len(interfaces) >= 1 {
				return nil
			}
			time.Sleep(5 * time.Second)
		}

		return fmt.Errorf("timed out waiting for NAT gateway %d to list an interface", gatewayID)
	}
}

func TestAccDataSourceNATGatewayInterfaces_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityLinodes,
		linodego.CapabilityVPCs,
	}, "core")
	if err != nil {
		t.Fatal(err)
	}

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the NAT gateway and interface, then wait for the
				// gateway to report the interface. The data sources may be
				// empty at this point; the next step re-reads them.
				Config: linodeinstancetmpl.ProviderNoPoll(t) +
					tmpl.DataBasic(t, label, region, "10.0.0.0/24"),
				Check: waitForNATGatewayInterface(natGatewayResName, 3*time.Minute),
			},
			{
				// Re-read the data sources now that the interface is listed.
				Config: linodeinstancetmpl.ProviderNoPoll(t) +
					tmpl.DataBasic(t, label, region, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("linode").AtMapKey("label"),
						knownvalue.StringExact(label),
					),
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("linode").AtMapKey("type"),
						knownvalue.StringExact("linode"),
					),
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("addresses"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						interfacesByGatewayRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("portsets"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceNATGatewayInterfaces_address(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityLinodes,
		linodego.CapabilityVPCs,
	}, "core")
	if err != nil {
		t.Fatal(err)
	}

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the manual NAT gateway and interface, then wait for the
				// gateway to report the interface.
				Config: linodeinstancetmpl.ProviderNoPoll(t) +
					tmpl.DataAddress(t, label, region, "10.0.0.0/24"),
				Check: waitForNATGatewayInterface(natGatewayResName, 3*time.Minute),
			},
			{
				// Re-read the address-scoped data source now that the interface is listed.
				Config: linodeinstancetmpl.ProviderNoPoll(t) +
					tmpl.DataAddress(t, label, region, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						interfacesByAddressRes,
						tfjsonpath.New("interfaces"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						interfacesByAddressRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("linode").AtMapKey("label"),
						knownvalue.StringExact(label),
					),
					statecheck.ExpectKnownValue(
						interfacesByAddressRes,
						tfjsonpath.New("interfaces").AtSliceIndex(0).AtMapKey("addresses"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
