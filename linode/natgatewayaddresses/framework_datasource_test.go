//go:build integration || natgatewayaddresses

package natgatewayaddresses_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgatewayaddresses/tmpl"
)

const (
	natAddressesDataSource = "data.linode_natgateway_addresses.foobar"
	natAddressesReserved   = "linode_networking_ip.reserved"
)

func TestAccDataSourceNATGatewayAddresses_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The label must begin with a letter.
	label := "tf-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("natgateway_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses").AtSliceIndex(0).AtMapKey("in_use"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses").AtSliceIndex(0).AtMapKey("interface_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses").AtSliceIndex(0).AtMapKey("interface_url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses").AtSliceIndex(0).AtMapKey("portset_assignments"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressesDataSource,
						tfjsonpath.New("addresses").AtSliceIndex(0).AtMapKey("portset_capacity"), knownvalue.NotNull()),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(natAddressesDataSource, "addresses.0.address", natAddressesReserved, "address"),
				),
			},
		},
	})
}
