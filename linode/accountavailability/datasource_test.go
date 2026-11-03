//go:build integration || accountavailability || act_tests

package accountavailability_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/accountavailability/tmpl"
)

func TestAccDataSourceNodeBalancers_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_availability.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(err)
	}
	region, err := acceptance.NewRegionSelector(context.Background(), client).
		WithCapabilities(linodego.CapabilityLinodes).
		WithSiteType("core").
		AvailableForAccount(context.Background(), linodego.CapabilityLinodes).
		Random()
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, region),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttrSet(resourceName, "unavailable.#"),
					resource.TestCheckResourceAttrSet(resourceName, "available.#"),
				),
			},
		},
	})
}
