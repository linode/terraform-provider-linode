//go:build integration || image

package image_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/image/tmpl"
)

func TestAccDataSourceImage_basic(t *testing.T) {
	t.Parallel()

	imageID := "linode/debian8"
	resourceName := "data.linode_image.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, imageID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", imageID),
					resource.TestCheckResourceAttr(resourceName, "label", "Debian 8"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", "manual"),
					resource.TestCheckResourceAttr(resourceName, "size", "1300"),
					resource.TestCheckResourceAttr(resourceName, "vendor", "Debian"),
					resource.TestCheckResourceAttrSet(resourceName, "capabilities.#"),
				),
			},
		},
	})
}

func TestAccDataSourceImage_replicate(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_image.foobar"
	imageName := acctest.RandomWithPrefix("tf_test")
	// TODO: Use random region once image gen2 works globally or with specific capabilities
	replicateRegion := "eu-west"

	file, err := createTempFile("tf-test-image-data-replicate-file", testImageBytes)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,

		CheckDestroy: checkImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataReplicate(t, imageName, file.Name(), testRegion, replicateRegion),
				Check: resource.ComposeTestCheckFunc(
					checkImageExists(resourceName, nil),
					resource.TestCheckResourceAttr(resourceName, "label", imageName),
					resource.TestCheckResourceAttr(resourceName, "description", "really descriptive text"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "size"),
					resource.TestCheckResourceAttr(resourceName, "type", "manual"),
					resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "deprecated"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(resourceName, "total_size"),
					resource.TestCheckResourceAttr(resourceName, "replications.#", "2"),
				),
			},
		},
	})
}
