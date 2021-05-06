package linode

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

// testImageBytes is a minimal Gzipped image.
// This is necessary because the API will reject invalid images.
var testImageBytes = []byte{0x1f, 0x8b, 0x08, 0x08, 0xbd, 0x5c, 0x91, 0x60,
	0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x69, 0x6d, 0x67, 0x00, 0x03, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

var testImageBytesNew = []byte{0x1f, 0x8b, 0x08, 0x08, 0x53, 0x13, 0x94, 0x60,
	0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x69, 0x6d, 0x67, 0x00, 0xcb, 0xc8,
	0xe4, 0x02, 0x00, 0x7a, 0x7a, 0x6f, 0xed, 0x03, 0x00, 0x00, 0x00}

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
					testAccCheckLinodeImageExists(resName, nil),
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
					testAccCheckLinodeImageExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "descriptive text"),
				),
			},
			{
				Config: testAccCheckLinodeImageConfigUpdates(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists(resName, nil),
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

func TestAccLinodeImage_uploadFile(t *testing.T) {
	t.Parallel()

	resName := "linode_image.foobar"
	imageName := acctest.RandomWithPrefix("tf_test")

	file, err := testAccCreateTempFile("tf-test-image-upload-file", testImageBytes)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var image linodego.Image

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeImageConfigUpload(imageName, file.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "really descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
					resource.TestCheckResourceAttrSet(resName, "file_hash"),
					resource.TestCheckResourceAttr(resName, "status", string(linodego.ImageStatusAvailable)),
				),
			},
			{
				PreConfig: func() {
					file.Write(testImageBytesNew)
				},
				Config: testAccCheckLinodeImageConfigUpload(imageName, file.Name()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "status", string(linodego.ImageStatusAvailable)),
				),
			},
		},
	})
}

func testAccCheckLinodeImageExists(name string, image *linodego.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		found, err := client.GetImage(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve state of image %s: %s", rs.Primary.Attributes["label"], err)
		}

		if image != nil {
			*image = *found
		}

		return nil
	}
}

func testAccCheckLinodeImageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client
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

func testAccCreateTempFile(name string, content []byte) (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %s", err)
	}

	if _, err := file.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %s", err)
	}

	return file, nil
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

func testAccCheckLinodeImageConfigUpload(image string, file string) string {
	return fmt.Sprintf(`
resource "linode_image" "foobar" {
	label = "%s"
	file_path = "%s"
	file_hash = filemd5("%s")
	region = "us-southeast"
	description = "really descriptive text"
}`, image, file, file)
}
