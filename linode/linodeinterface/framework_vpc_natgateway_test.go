//go:build integration || linodeinterface

package linodeinterface_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	linodeinstancetmpl "github.com/linode/terraform-provider-linode/v4/linode/acceptance/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/linodeinterface/tmpl"
)

func TestAccLinodeInterface_vpc_natgateway(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")

	natGatewayRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityLinodes,
		linodego.CapabilityVPCs,
	}, "core")
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	natGatewayPath := func(key string) tfjsonpath.Path {
		return tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("natgateway").AtMapKey(key)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCNATGateway(t, label, natGatewayRegion, "10.0.0.0/24"),
			},
			{
				// Refresh step: read the interface after it has been created on the NAT gateway subnet.
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.VPCNATGateway(t, label, natGatewayRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						testInterfaceResName, natGatewayPath("id"), knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName, natGatewayPath("label"), knownvalue.StringExact(label+"-natgateway"),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName, natGatewayPath("addresses"), knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName, natGatewayPath("portset_assignments"), knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						testInterfaceResName, natGatewayPath("portset_capacity"), knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
