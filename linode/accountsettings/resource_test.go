//go:build integration || accountsettings || act_tests

package accountsettings_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("backups_enabled"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("managed"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("network_helper"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("object_storage"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("interfaces_for_new_linodes"), knownvalue.NotNull()),
				},
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
	currInterfacesForNewLinodes := accountSettings.InterfacesForNewLinodes
	currMaintenancePolicy := accountSettings.MaintenancePolicy

	updatedLongviewPlan := "longview-10"
	updatedBackupsEnabled := !currBackupsEnabled
	updatedNetworkHelper := !currNetworkHelper
	updatedMaintenancePolicy := "linode/power_off_on"

	var updatedInterfacesForNewLinodes string
	if currInterfacesForNewLinodes == linodego.LegacyConfigDefaultButLinodeAllowed {
		updatedInterfacesForNewLinodes = string(linodego.LinodeDefaultButLegacyConfigAllowed)
	} else {
		updatedInterfacesForNewLinodes = string(linodego.LegacyConfigDefaultButLinodeAllowed)
	}

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
				Config: tmpl.Updates(t, updatedLongviewPlan, updatedInterfacesForNewLinodes, updatedBackupsEnabled, updatedNetworkHelper, updatedMaintenancePolicy),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("longview_subscription"), knownvalue.StringExact(updatedLongviewPlan)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("maintenance_policy"), knownvalue.StringExact(updatedMaintenancePolicy)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("backups_enabled"), knownvalue.Bool(updatedBackupsEnabled)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("network_helper"), knownvalue.Bool(updatedNetworkHelper)),
					statecheck.ExpectKnownValue(
						resourceName, tfjsonpath.New("interfaces_for_new_linodes"), knownvalue.StringExact(updatedInterfacesForNewLinodes),
					),
				},
			},
			{
				Config: tmpl.Updates(t, currLongviewPlan, string(currInterfacesForNewLinodes), currBackupsEnabled, currNetworkHelper, currMaintenancePolicy),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("longview_subscription"), knownvalue.StringExact(currLongviewPlan)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("maintenance_policy"), knownvalue.StringExact(currMaintenancePolicy)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("backups_enabled"), knownvalue.Bool(currBackupsEnabled)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("network_helper"), knownvalue.Bool(currNetworkHelper)),
					statecheck.ExpectKnownValue(
						resourceName, tfjsonpath.New("interfaces_for_new_linodes"), knownvalue.StringExact(string(currInterfacesForNewLinodes)),
					),
				},
			},
		},
	})
}
