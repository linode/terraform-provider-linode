//go:build integration || kernal

package kernel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/kernel/tmpl"
)

func TestAccDataSourceKernel_basic(t *testing.T) {
	t.Parallel()

	kernelID := "linode/latest-64bit"
	resourceName := "data.linode_kernel.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
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
					resource.TestCheckResourceAttrSet(resourceName, "built"),
				),
			},
		},
	})
}
