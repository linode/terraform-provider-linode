//go:build integration || networkreservedips

package networkreservedips_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkreservedips/tmpl"
)

func TestAccDataSource_reservedIPList(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_reserved_ip_list.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataList(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "reserved_ips.#"),
				),
			},
		},
	})
}
