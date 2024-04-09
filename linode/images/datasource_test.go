//go:build integration || images

package images_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/images/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil)
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceImages_basic_smoke(t *testing.T) {
	t.Parallel()

	imageName := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_images.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, imageName, testRegion),
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
					resource.TestCheckResourceAttrSet(resourceName, "images.0.capabilities.#"),
					resource.TestCheckResourceAttr(resourceName, "images.1.label", imageName),
					resource.TestCheckResourceAttr(resourceName, "images.1.description", "descriptive text"),
					resource.TestCheckResourceAttr(resourceName, "images.1.is_public", "false"),
					resource.TestCheckResourceAttr(resourceName, "images.1.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.created"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.size"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.deprecated"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.capabilities.#"),
				),
			},

			// These cases are all used in the same test to avoid recreating images unnecessarily
			{
				Config: tmpl.DataLatest(t, imageName, testRegion),
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
				Config: tmpl.DataLatestEmpty(t, imageName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "images.#", "0"),
				),
			},

			{
				Config: tmpl.DataOrder(t, imageName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					// Ensure order is correctly appended to filter
					resource.TestCheckResourceAttr(resourceName, "images.#", "2"),
				),
			},

			{
				Config: tmpl.DataSubstring(t, imageName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					// Ensure order is correctly appended to filter
					acceptance.CheckResourceAttrGreaterThan(resourceName, "images.#", 1),
					acceptance.CheckResourceAttrContains(resourceName, "images.0.label", "Alpine"),
				),
			},

			{
				Config: tmpl.DataClientFilter(t, imageName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "images.#", "1"),
					acceptance.CheckResourceAttrContains(resourceName, "images.0.label", imageName),
				),
			},
		},
	})
}
