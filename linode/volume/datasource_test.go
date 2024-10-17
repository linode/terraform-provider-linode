//go:build integration || volume

package volume_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/volume/tmpl"
)

func TestAccDataSourceVolume_basic(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", testRegion),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "label", volumeName),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resourceName, "linode_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "disabled"),
				),
			},
		},
	})
}

func TestAccDataSourceVolume_withBlockStorageEncryption(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_volume.foobar"

	// Resolve a region with support for Block Storage Encryption
	targetRegion, err := acceptance.GetRandomRegionWithCaps(
		[]string{"Linodes", "Block Storage Encryption"},
		"core",
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithBlockStorageEncryption(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", targetRegion),
					resource.TestCheckResourceAttr(resourceName, "encryption", "enabled"),
				),
			},
		},
	})
}
