//go:build integration || instancenetworking

package instancenetworking_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking/tmpl"
)

const testInstanceNetworkResName = "data.linode_instance_networking.test"

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityVPCs, linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceInstanceNetworking_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, name, testRegion),
			},
			{
				Config: tmpl.DataBasic(t, name, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.private.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.public.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.reserved.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.shared.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.global.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.link_local.%"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.slaac.%"),
				),
			},
		},
	})
}

func TestAccDataSourceInstanceNetworking_vpc(t *testing.T) {
	t.Parallel()

	instanceVPCIP := "10.0.0.3"
	name := acctest.RandomWithPrefix("tf-test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPC(t, name, testRegion, "10.0.0.0/24", instanceVPCIP),
			},
			{
				Config: tmpl.DataVPC(t, name, testRegion, "10.0.0.0/24", instanceVPCIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.vpc.#"),
					resource.TestCheckResourceAttr(testInstanceNetworkResName, "ipv4.0.vpc.0.address", instanceVPCIP),
				),
			},
		},
	})
}

func TestAccDataSourceInstanceNetworking_basicwithReseved(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic_withReservedField(t, name, testRegion),
			},
			{
				Config: tmpl.DataBasic_withReservedField(t, name, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.private.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.public.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.reserved.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.shared.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.global.#"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.link_local.%"),
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv6.0.slaac.%"),
					// The reserved IP appears in the public list; verify 2 public IPs are present
					resource.TestCheckResourceAttr(testInstanceNetworkResName, "ipv4.0.public.#", "2"),
					// At least one public IP has assigned_entity populated (the reserved IP)
					resource.TestCheckResourceAttrSet(testInstanceNetworkResName, "ipv4.0.public.0.assigned_entity.type"),
				),
			},
		},
	})
}

// TODO (Linode Interfaces): Add test for new interface_id field.

func TestAccDataSourceInstanceNetworking_vpcDualStack(t *testing.T) {
	t.Parallel()

	instanceVPCIP := "10.0.0.5"
	name := acctest.RandomWithPrefix("tf-test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPCDualStack(t, name, testRegion, "10.0.0.0/24", instanceVPCIP),
				ConfigStateChecks: []statecheck.StateCheck{
					// IPv4 VPC IPs are present
					statecheck.ExpectKnownValue(testInstanceNetworkResName,
						tfjsonpath.New("ipv4").AtSliceIndex(0).AtMapKey("vpc"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInstanceNetworkResName,
						tfjsonpath.New("ipv4").AtSliceIndex(0).AtMapKey("vpc").AtSliceIndex(0).AtMapKey("address"),
						knownvalue.StringExact(instanceVPCIP)),
					// IPv6 VPC IPs are present
					statecheck.ExpectKnownValue(testInstanceNetworkResName,
						tfjsonpath.New("ipv6").AtSliceIndex(0).AtMapKey("vpc"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInstanceNetworkResName,
						tfjsonpath.New("ipv6").AtSliceIndex(0).AtMapKey("vpc").AtSliceIndex(0).AtMapKey("ipv6_range"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testInstanceNetworkResName,
						tfjsonpath.New("ipv6").AtSliceIndex(0).AtMapKey("vpc").AtSliceIndex(0).AtMapKey("ipv6_addresses"),
						knownvalue.NotNull()),
				},
			},
		},
	})
}
