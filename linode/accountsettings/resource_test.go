package accountsettings_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/accountsettings/tmpl"
)

func TestAccResourceAccountSettings_basic(t *testing.T) {
	acceptance.OptInTest(t)

	resourceName := "linode_account_settings.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t),
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

func TestAccResourceAccountSettings_update(t *testing.T) {
	acceptance.OptInTest(t)

	resourceName := "linode_account_settings.foobar"

	accountSettings := linodego.AccountSettings{}
	longviewSettings := linodego.LongviewPlan{}

	currLongviewPlan := longviewSettings.ID
	currBackupsEnabled := accountSettings.BackupsEnabled
	currNetworkHelper := accountSettings.NetworkHelper

	updatedLongviewPlan := "longview-10"
	updatedBackupsEnabled := !currBackupsEnabled
	updatedNetworkHelper := !currNetworkHelper

	if currLongviewPlan == "" || currLongviewPlan == "longview-10" {
		updatedLongviewPlan = "longview-3"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Updates(t, updatedLongviewPlan, updatedBackupsEnabled, updatedNetworkHelper),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "longview_subscription", updatedLongviewPlan),
					resource.TestCheckResourceAttr(resourceName, "backups_enabled", strconv.FormatBool(updatedBackupsEnabled)),
					resource.TestCheckResourceAttr(resourceName, "network_helper", strconv.FormatBool(updatedNetworkHelper)),
				),
			},
			{
				Config: tmpl.Updates(t, currLongviewPlan, currBackupsEnabled, currNetworkHelper),
			},
		},
	})
}
