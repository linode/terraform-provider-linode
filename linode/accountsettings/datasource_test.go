package accountsettings_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/accountsettings/tmpl"
)

func TestAccDataSourceLinodeAccountSettings_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_settings.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backups_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "managed"),
					resource.TestCheckResourceAttrSet(resourceName, "network_helper"),
					resource.TestCheckResourceAttrSet(resourceName, "object_storage"),
				),
			},
		},
	})
}
