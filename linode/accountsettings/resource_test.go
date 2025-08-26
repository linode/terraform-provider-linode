//go:build integration || accountsettings || act_tests

package accountsettings_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/accountsettings/tmpl"
)

func TestAccResourceAccountSettings_basic(t *testing.T) {
	acceptance.OptInTest(t)

	resourceName := "linode_account_settings.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
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

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	accountSettings, _ := client.GetAccountSettings(context.Background())
	longviewSettings, _ := client.GetLongviewPlan(context.Background())

	currLongviewPlan := longviewSettings.ID
	currBackupsEnabled := accountSettings.BackupsEnabled
	currNetworkHelper := accountSettings.NetworkHelper
	currMaintenancePolicy := accountSettings.MaintenancePolicy

	updatedLongviewPlan := "longview-10"
	updatedBackupsEnabled := !currBackupsEnabled
	updatedNetworkHelper := !currNetworkHelper
	updatedMaintenancePolicy := "linode/power_off_on"

	if currLongviewPlan == "" || currLongviewPlan == "longview-10" {
		updatedLongviewPlan = "longview-3"
	}

	if currMaintenancePolicy == "" || currMaintenancePolicy == "linode/power_off_on" {
		updatedLongviewPlan = "linode/migrate"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Updates(t, updatedLongviewPlan, updatedBackupsEnabled, updatedNetworkHelper, updatedMaintenancePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "longview_subscription", updatedLongviewPlan),
					resource.TestCheckResourceAttr(resourceName, "backups_enabled", strconv.FormatBool(updatedBackupsEnabled)),
					resource.TestCheckResourceAttr(resourceName, "network_helper", strconv.FormatBool(updatedNetworkHelper)),
					resource.TestCheckResourceAttr(resourceName, "maintenance_policy", updatedMaintenancePolicy),
				),
			},
			{
				Config: tmpl.Updates(t, currLongviewPlan, currBackupsEnabled, currNetworkHelper, currMaintenancePolicy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "longview_subscription", currLongviewPlan),
					resource.TestCheckResourceAttr(resourceName, "backups_enabled", strconv.FormatBool(currBackupsEnabled)),
					resource.TestCheckResourceAttr(resourceName, "network_helper", strconv.FormatBool(currNetworkHelper)),
					resource.TestCheckResourceAttr(resourceName, "maintenance_policy", currMaintenancePolicy),
				),
			},
		},
	})
}
