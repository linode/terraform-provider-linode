//go:build integration || natgatewayaddress

package natgatewayaddress_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgatewayaddress/tmpl"
)

const (
	natAddressDataSource = "data.linode_natgateway_address.foobar"
	natAddressReserved   = "linode_networking_ip.reserved"
)

func TestAccDataSourceNATGatewayAddress_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Fatal(err)
	}

	label := "tf-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("natgateway_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("in_use"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("interface_count"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("interface_url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("portset_assignments"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natAddressDataSource,
						tfjsonpath.New("portset_capacity"), knownvalue.NotNull()),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(natAddressDataSource, "address", natAddressReserved, "address"),
				),
			},
		},
	})
}
