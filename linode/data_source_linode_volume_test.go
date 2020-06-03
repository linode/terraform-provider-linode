package linode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeVolume_basic(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	resourceName := "data.linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName) + testDataSourceLinodeVolumeByID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "region", "us-west"),
					resource.TestCheckResourceAttr(resourceName, "size", "20"),
					resource.TestCheckResourceAttr(resourceName, "label", volumeName),
					resource.TestCheckResourceAttr(resourceName, "tags.4106436895", "tf_test"),
					resource.TestCheckResourceAttr(resourceName, "linode_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
				),
			},
		},
	})
}

func testDataSourceLinodeVolumeByID() string {
	return `
data "linode_volume" "foobar" {
	id = "${linode_volume.foobar.id}"
}`
}
