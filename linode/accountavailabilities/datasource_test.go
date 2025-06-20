//go:build integration || accountavailabilities || act_tests

package accountavailabilities_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailabilities/tmpl"
)

func TestAccDataSourceAccountAvailabilities_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_availabilities.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "availabilities.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "availabilities.0.unavailable.#"),
					resource.TestCheckResourceAttrSet(resourceName, "availabilities.0.available.#"),
				),
			},
			{
				Config: tmpl.DataFilterRegion(t, "us-east"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availabilities.0.region", "us-east"),
				),
			},
		},
	})
}
