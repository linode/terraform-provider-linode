package images_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func TestAccDataSourceImages_basic(t *testing.T) {
	t.Parallel()

	imageName := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_images.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigBasic(imageName),
				Check: resource.ComposeTestCheckFunc(
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
		},
	})
}

func imageConfigBasic(image string) string {
	return fmt.Sprintf(`
	resource "linode_instance" "foobar" {
		label = "%s"
		group = "tf_test"
		type = "g6-standard-1"
		region = "us-east"
		disk {
			label = "disk"
			size = 1000
			filesystem = "ext4"
		}
	}
	resource "linode_image" "foobar" {
		linode_id = "${linode_instance.foobar.id}"
		disk_id = "${linode_instance.foobar.disk.0.id}"
		label = "%s"
		description = "descriptive text"
	}`, image, image)
}

func dataSourceConfigBasic(image string) string {
	return imageConfigBasic(image) + `
data "linode_images" "foobar" {
	filter {
		name = "label"
		values = [linode_image.foobar.label]
	}

	filter {
		name = "is_public"
		values = ["false"]
	}
}`
}
