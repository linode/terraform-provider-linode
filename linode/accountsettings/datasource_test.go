//go:build integration || accountsettings

package accountsettings_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/accountsettings/tmpl"
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "backups_enabled", strconv.FormatBool(settings.BackupsEnabled)),
					resource.TestCheckResourceAttr(resourceName, "managed", strconv.FormatBool(settings.Managed)),
					resource.TestCheckResourceAttr(resourceName, "network_helper", strconv.FormatBool(settings.NetworkHelper)),
					resource.TestCheckResourceAttr(resourceName, "object_storage", objectStorageVal),
					resource.TestCheckResourceAttr(resourceName, "longview_subscription", longviewVal),
				),
			},
		},
	})
}
