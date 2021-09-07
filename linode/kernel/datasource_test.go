package kernel_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/kernel/tmpl"

	"testing"
)

func TestAccDataSourceKernel_basic(t *testing.T) {
	t.Parallel()

	kernelID := "linode/latest-64bit"
	resourceName := "data.linode_kernel.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, kernelID),
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
