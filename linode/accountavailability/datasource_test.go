//go:build integration || accountavailability

package accountavailability_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailability/tmpl"
)

func TestAccDataSourceNodeBalancers_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_availability.foobar"

	region, err := acceptance.GetRandomRegionWithCaps(nil)
	if err != nil {
		log.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, region),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttrSet(resourceName, "unavailable.#"),
				),
			},
		},
	})
}
