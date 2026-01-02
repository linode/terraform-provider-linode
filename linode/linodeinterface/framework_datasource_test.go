//go:build integration || linodeinterface

package linodeinterface_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/linodeinterface/tmpl"
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
