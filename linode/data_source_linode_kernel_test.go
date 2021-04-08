package linode

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLinodeKernel_basic(t *testing.T) {
	t.Parallel()

	kernelID := "linode/latest-64bit"
	resourceName := "data.linode_kernel.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeKernelBasic(kernelID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", kernelID),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "architecture"),
					resource.TestCheckResourceAttrSet(resourceName, "deprecated"),
					resource.TestCheckResourceAttrSet(resourceName, "kvm"),
					resource.TestCheckResourceAttrSet(resourceName, "label"),
					resource.TestCheckResourceAttrSet(resourceName, "pvops"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "xen"),
				),
			},
		},
	})
}

func testDataSourceLinodeKernelBasic(kernelID string) string {
	return fmt.Sprintf(`
data "linode_kernel" "foobar" {
	id = "%s"
}`, kernelID)
}
