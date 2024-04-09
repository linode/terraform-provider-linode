//go:build integration || kernels

package kernels_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/kernels/tmpl"
)

func TestAccDataSourceKernels_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_kernels.kernels"

	kernelID := "linode/latest-64bit"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, kernelID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kernels.0.id", kernelID),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.architecture"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.deprecated"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.kvm"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.pvops"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.version"),
					resource.TestCheckResourceAttrSet(resourceName, "kernels.0.xen"),
				),
			},
			{
				Config: tmpl.DataFilter(t, kernelID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kernels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.id", kernelID),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.deprecated", "false"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.kvm", "true"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.label", "Latest 64 bit (6.2.9-x86_64-linode160)"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.pvops", "true"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.version", "6.2.9"),
					resource.TestCheckResourceAttr(resourceName, "kernels.0.xen", "false"),
				),
			},
			{
				Config: tmpl.DataFilterEmpty(t, kernelID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "kernels.#", "0"),
				),
			},
		},
	})
}
