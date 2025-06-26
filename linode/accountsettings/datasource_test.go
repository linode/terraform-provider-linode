//go:build integration || accountsettings || act_tests

package accountsettings_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/accountsettings/tmpl"
)

func TestAccDataSourceLinodeAccountSettings_basic(t *testing.T) {
	acceptance.OptInTest(t)

	t.Parallel()

	resourceName := "data.linode_account_settings.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	settings, err := client.GetAccountSettings(context.Background())
	if err != nil {
		t.Fatalf("failed to get account settings: %s", err)
	}

	objectStorageVal := ""
	if settings.ObjectStorage != nil {
		objectStorageVal = *settings.ObjectStorage
	}

	longviewVal := ""
	if settings.LongviewSubscription != nil {
		longviewVal = *settings.LongviewSubscription
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("backups_enabled"), knownvalue.Bool(settings.BackupsEnabled)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("managed"), knownvalue.Bool(settings.Managed)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("network_helper"), knownvalue.Bool(settings.NetworkHelper)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("object_storage"), knownvalue.StringExact(objectStorageVal)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("longview_subscription"), knownvalue.StringExact(longviewVal)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("interfaces_for_new_linodes"), knownvalue.StringExact(string(settings.InterfacesForNewLinodes))),
				},
			},
		},
	})
}
