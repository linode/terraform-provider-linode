//go:build integration || image

package image_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/image/tmpl"
)

// testImageBytes is a minimal Gzipped image.
// This is necessary because the API will reject invalid images.
var testImageBytes = []byte{
	0x1f, 0x8b, 0x08, 0x08, 0xbd, 0x5c, 0x91, 0x60,
	0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x69, 0x6d, 0x67, 0x00, 0x03, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

// MD5 digest of testImageBytes
const testImageMD5 = "0cf442194905e7be019a11660df8164f"

var testImageBytesNew = []byte{
	0x1f, 0x8b, 0x08, 0x08, 0x53, 0x13, 0x94, 0x60,
	0x00, 0x03, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x69, 0x6d, 0x67, 0x00, 0xcb, 0xc8,
	0xe4, 0x02, 0x00, 0x7a, 0x7a, 0x6f, 0xed, 0x03, 0x00, 0x00, 0x00,
}

var (
	testRegion  string
	testRegions []string
)

func init() {
	resource.AddTestSweepers("linode_image", &resource.Sweeper{
		Name: "linode_image",
		F:    sweep,
	})

	regions, err := acceptance.GetRegionsWithCaps([]string{"Object Storage"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = regions[1]
	testRegions = regions
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	images, err := client.ListImages(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting images: %s", err)
	}
	for _, image := range images {
		if !acceptance.ShouldSweep(prefix, image.Label) {
			continue
		}
		err := client.DeleteImage(context.Background(), image.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", image.Label, err)
		}
	}

	return nil
}

func TestAccImage_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_image.foobar"
	imageName := acctest.RandomWithPrefix("tf_test")
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,

		CheckDestroy: checkImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, imageName, testRegion, label, "test-tag"),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttr(resName, "capabilities.0", "cloud-init"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
				),
			},
			{
				ResourceName: resName,
				ImportState:  true,
				// ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id", "disk_id", "firewall_id"},
			},
		},
	})
}

func TestAccImage_update(t *testing.T) {
	t.Parallel()

	imageName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_image.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkImageDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, imageName, testRegion, label, "test-tag"),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "capabilities.#"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "test-tag"),
				),
			},
			{
				Config: tmpl.Updates(t, imageName, testRegion, label, "updated-tag"),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", imageName)),
					resource.TestCheckResourceAttr(resName, "description", "more descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "updated-tag"),
				),
			},
			{
				ResourceName: resName,
				ImportState:  true,
				// ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id", "disk_id", "firewall_id"},
			},
		},
	})
}

func TestAccImage_uploadFile(t *testing.T) {
	t.Parallel()

	resName := "linode_image.foobar"
	imageName := acctest.RandomWithPrefix("tf_test")

	file, err := createTempFile("tf-test-image-upload-file", testImageBytes)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var image linodego.Image

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Upload(t, imageName, file.Name(), testRegion, "test-tag"),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "description", "really descriptive text"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "created_by"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "type", "manual"),
					resource.TestCheckResourceAttr(resName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resName, "deprecated"),
					resource.TestCheckResourceAttr(resName, "file_hash", testImageMD5),
					resource.TestCheckResourceAttr(resName, "status", string(linodego.ImageStatusAvailable)),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
				),
			},
			{
				PreConfig: func() {
					file.Write(testImageBytesNew)
				},
				Config: tmpl.Upload(t, imageName, file.Name(), testRegion, "test-tag"),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "status", string(linodego.ImageStatusAvailable)),
				),
			},
		},
	})
}

func TestAccImage_replicate(t *testing.T) {
	t.Parallel()

	resName := "linode_image.foobar"
	imageName := acctest.RandomWithPrefix("tf_test")

	if len(testRegions) < 4 {
		t.Skipf("Not enough number of capable regions for image replication test. Skipping now...")
	}

	file, err := createTempFile("tf-test-image-replicate-file", testImageBytes)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	var image linodego.Image

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Replicate(t, imageName, file.Name(), testRegion, testRegions[0]),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "replications.#", "2"),
				),
			},
			{
				// Remove the one of the available region and replicate the image in a new region
				Config: tmpl.Replicate(t, imageName, file.Name(), testRegion, testRegions[2]),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resName, &image),
					resource.TestCheckResourceAttr(resName, "label", imageName),
					resource.TestCheckResourceAttr(resName, "replications.#", "2"),
				),
			},
			{
				// Remove all available region and replicate the image in new regions
				Config: tmpl.Replicate(t, imageName, file.Name(), testRegions[0], testRegions[3]),
				ExpectError: regexp.MustCompile(
					"At least one available region must be specified"),
			},
			{
				// Remove all available region
				Config: tmpl.NoReplicaRegions(t, imageName, file.Name(), testRegion),
				ExpectError: regexp.MustCompile(
					"At least one available region must be specified"),
			},
		},
	})
}

func checkImageExists(name string, image *linodego.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkImageDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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

func createTempFile(name string, content []byte) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), name)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %s", err)
	}

	if _, err := file.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %s", err)
	}

	return file, nil
}
