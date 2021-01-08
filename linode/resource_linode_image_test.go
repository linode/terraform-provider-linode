package linode

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_image", &resource.Sweeper{
		Name: "linode_image",
		F:    testSweepLinodeImage,
	})

}

func testSweepLinodeImage(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	images, err := client.ListImages(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting images: %s", err)
	}
	for _, image := range images {
		if !shouldSweepAcceptanceTestResource(prefix, image.Label) {
			continue
		}
		err := client.DeleteImage(context.Background(), image.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", image.Label, err)
		}
	}

	return nil
}

func TestAccLinodeImage_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_image.foobar"
	var ImageName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeImageConfigBasic(ImageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists,
					resource.TestCheckResourceAttr(resName, "label", ImageName),
					resource.TestCheckResourceAttr(resName, "description", "descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id", "disk_id"},
			},
		},
	})
}

func TestAccLinodeImage_update(t *testing.T) {
	t.Parallel()

	var imageName = acctest.RandomWithPrefix("tf_test")
	var resName = "linode_image.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeImageConfigBasic(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists,
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "descriptive text"),
				),
			},
			{
				Config: testAccCheckLinodeImageConfigUpdates(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", imageName)),
					resource.TestCheckResourceAttr(resName, "description", "more descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id", "disk_id"},
			},
		},
	})
}

func testAccCheckLinodeImageExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_Image" {
			continue
		}

		if _, err := client.GetImage(context.Background(), rs.Primary.ID); err != nil {
			return fmt.Errorf("Error retrieving state of Image %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeImageDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_Image" {
			continue
		}

		_, err := client.GetImage(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Linode Image with id %s still exists", rs.Primary.ID)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Image with id %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLinodeImageConfigBasic(image string) string {
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

func testAccCheckLinodeImageConfigUpdates(image string) string {
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
		label = "%s_renamed"
		description = "more descriptive text"
	}`, image, image)
}
