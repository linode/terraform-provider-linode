//go:build integration || accountavailabilities

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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "availabilities.0.region"),
					resource.TestCheckResourceAttrSet(resourceName, "availabilities.0.unavailable.#"),
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
