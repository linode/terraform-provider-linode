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

func TestAccDataSourceLinodeInterface_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_interface.test"
	linodeLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataSourceBasic(t, linodeLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "linode_id"),
					resource.TestCheckResourceAttrSet(resourceName, "default_route.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "default_route.ipv6"),
					resource.TestCheckResourceAttrSet(resourceName, "public.ipv4.addresses.#"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeInterface_vpc_natgateway(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf-test")
	dataSourceName := "data.linode_interface.test"

	natGatewayRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityLinodes,
		linodego.CapabilityVPCs,
	}, "core")
	if err != nil {
		log.Fatal(fmt.Errorf("Error getting region: %s", err))
	}

	// Build a fresh path per check: tfjsonpath.Path.AtMapKey appends to a shared slice.
	natGatewayPath := func(key string) tfjsonpath.Path {
		return tfjsonpath.New("vpc").AtMapKey("ipv4").AtMapKey("natgateway").AtMapKey(key)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: linodeinstancetmpl.ProviderNoPoll(t) + tmpl.DataSourceVPCNATGateway(t, label, natGatewayRegion, "10.0.0.0/24"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("linode_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("vpc").AtMapKey("subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, natGatewayPath("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, natGatewayPath("label"), knownvalue.StringExact(label+"-natgateway")),
					statecheck.ExpectKnownValue(dataSourceName, natGatewayPath("addresses"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, natGatewayPath("portset_assignments"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, natGatewayPath("portset_capacity"), knownvalue.NotNull()),
				},
			},
		},
	})
}
