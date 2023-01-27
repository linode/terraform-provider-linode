package lkeversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/lkeversion/tmpl"
)

func TestAccDataSourceLinodeLkeVersion_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_version.foobar"

	lkeVersionID := "1.25"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, lkeVersionID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", lkeVersionID),
				),
			},
		},
	})
}
