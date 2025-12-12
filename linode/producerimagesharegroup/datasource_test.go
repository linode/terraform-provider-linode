//go:build integration || producerimagesharegroup

package producerimagesharegroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroup/tmpl"
)

func TestAccDataSourceImageShareGroup_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_producer_image_share_group.foobar"

	label := acctest.RandomWithPrefix("tf-test")
	description := "A cool description."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", label),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckNoResourceAttr(resourceName, "updated"),
					resource.TestCheckNoResourceAttr(resourceName, "expiry"),
				),
			},
		},
	})
}
