//go:build integration || lketypes

package lketypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lketypes/tmpl"
)

func TestAccDataSourceLKETypes_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_lke_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "types.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.id", "lke-sa"),
					resource.TestCheckResourceAttr(dataSourceName, "types.0.label", "LKE Standard Availability"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.transfer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(dataSourceName, "types.0.price.0.monthly"),
				),
			},
		},
	})
}
