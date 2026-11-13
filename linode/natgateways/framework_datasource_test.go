//go:build integration || natgateways

package natgateways_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgateways/tmpl"
)

const (
	natGatewaysByLabel = "data.linode_natgateways.by_label"
	natGatewaysByID    = "data.linode_natgateways.by_id"
	natGatewayFirst    = "linode_natgateway.first"
	natGatewaySecond   = "linode_natgateway.second"
)

func TestAccDataSourceNATGateways_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityNATGateway}, "core")
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
					// Label filter returns only the first gateway.
					statecheck.ExpectKnownValue(
						natGatewaysByLabel,
						tfjsonpath.New("natgateways"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						natGatewaysByLabel,
						tfjsonpath.New("natgateways").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label+"-first"),
					),
					statecheck.ExpectKnownValue(
						natGatewaysByLabel,
						tfjsonpath.New("natgateways").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.StringExact(region),
					),

					// Id filter returns only the second gateway.
					statecheck.ExpectKnownValue(
						natGatewaysByID,
						tfjsonpath.New("natgateways"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						natGatewaysByID,
						tfjsonpath.New("natgateways").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label+"-second"),
					),
					statecheck.ExpectKnownValue(
						natGatewaysByID,
						tfjsonpath.New("natgateways").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.StringExact(region),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(natGatewaysByLabel, "natgateways.0.id", natGatewayFirst, "id"),
					resource.TestCheckResourceAttrPair(natGatewaysByID, "natgateways.0.id", natGatewaySecond, "id"),
				),
			},
		},
	})
}
