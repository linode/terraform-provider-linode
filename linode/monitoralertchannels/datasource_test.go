package monitoralertchannels_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertchannels/tmpl"
)

func TestAccDataSourceMonitorAlertChannels_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_monitor_alert_channels.channels"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("monitor_alert_channels"), knownvalue.NotNull()),

					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("channel_type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("created"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("updated"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("created_by"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("updated_by"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("alerts").AtMapKey("url"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("alerts").AtMapKey("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("alerts").AtMapKey("alert_count"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: tmpl.DataFilter(t, "system"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("monitor_alert_channels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("monitor_alert_channels").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.StringExact("system"),
					),
				},
			},
		},
	})
}
