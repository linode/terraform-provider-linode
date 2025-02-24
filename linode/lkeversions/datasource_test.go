//go:build integration || lkeversions

package lkeversions_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversions/tmpl"
)

func TestAccDataSourceLinodeLkeVersions_NoTier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_versions.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataNoTier(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "versions.0.id"),
					resource.TestCheckNoResourceAttr(resourceName, "versions.0.tier"),
					resource.TestCheckNoResourceAttr(resourceName, "tier"),
				),
			},
		},
	})
}

func TestAccDataSourceLinodeLkeVersions_Tier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_versions.foobar"
	tier := "enterprise"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataTier(t, tier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "versions.0.id"),
					resource.TestCheckResourceAttr(resourceName, "versions.0.tier", tier),
					resource.TestCheckResourceAttr(resourceName, "tier", tier),
				),
			},
		},
	})
}
