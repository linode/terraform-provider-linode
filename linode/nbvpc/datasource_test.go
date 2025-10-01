//go:build integration || nbvpc

package nbvpc_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/nbvpc/tmpl"
)

func TestAccDataSource_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nodebalancer_vpc.test"

	label := acctest.RandomWithPrefix("tf-test")

	targetRegion, err := acceptance.GetRandomRegionWithCaps([]string{"NodeBalancers", "VPCs"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, targetRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("nodebalancer_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("vpc_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("subnet_id"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("ipv4_range"),
						knownvalue.StringExact("10.0.0.4/30"),
					),
				},
			},
		},
	})
}
