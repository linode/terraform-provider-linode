package images_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/images/tmpl"
)

func TestAccDataSourceImages_basic(t *testing.T) {
	t.Parallel()

	imageName := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_images.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "images.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "images.0.label", imageName),
					resource.TestCheckResourceAttr(resourceName, "images.0.description", "descriptive text"),
					resource.TestCheckResourceAttr(resourceName, "images.0.is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "images.0.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.deprecated"),
					resource.TestCheckResourceAttr(resourceName, "images.1.label", imageName),
					resource.TestCheckResourceAttr(resourceName, "images.1.description", "descriptive text"),
					resource.TestCheckResourceAttr(resourceName, "images.1.is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "images.1.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.created"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.size"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.deprecated"),
				),
			},

			// These cases are all used in the same test to avoid recreating images unnecessarily
			{
				Config: tmpl.DataLatest(t, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "images.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "images.0.label", imageName),
					resource.TestCheckResourceAttr(resourceName, "images.0.description", "descriptive text"),
					resource.TestCheckResourceAttr(resourceName, "images.0.is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "images.0.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.deprecated"),
				),
			},

			{
				Config: tmpl.DataLatestEmpty(t, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "images.#", "0"),
				),
			},

			{
				Config: tmpl.DataOrder(t, imageName),
				Check: resource.ComposeTestCheckFunc(
					// Ensure order is correctly appended to filter
					resource.TestCheckResourceAttr(resourceName, "images.#", "2"),
					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order_by\":\"size\""),
					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order\":\"desc\""),
				),
			},
		},
	})
}
