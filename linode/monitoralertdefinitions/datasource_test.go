//go:build integration || monitoralertdefinitions

package monitoralertdefinitions_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinitions/tmpl"
)

func TestAccDataSourceAlertDefinitions_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_monitor_alert_definitions.foobar"
	alertLabel := acctest.RandomWithPrefix("tf-test")

	// TODO: revert to use alert channels from API once it's available
	//client, err := acceptance.GetTestClient()
	//if err != nil {
	//	fmt.Errorf("Error getting client: %s", err)
	//}
	//
	//channels, err := client.ListAlertChannels(context.Background(), nil)
	//if err != nil {
	//	fmt.Errorf("error listing alert channels: %s", err)
	//}
	//
	//if len(channels) < 1 {
	//	t.Skipf("Skipping test: At least one alert channel is required for alert definition tests")
	//}
	//
	//channelID := channels[0].ID
	channelID := 10000

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, alertLabel, channelID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("service_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("channel_ids"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("description"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("entity_ids"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("severity"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("rule_criteria"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("has_more_resources"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("alert_channels"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("created_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("updated_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("class"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("alert_channels").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("alert_channels").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("alert_channels").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("alert_definitions").AtSliceIndex(0).AtMapKey("alert_channels").AtSliceIndex(0).AtMapKey("url"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceAlertDefinitions_filter(t *testing.T) {
	t.Parallel()

	resName := "data.linode_monitor_alert_definitions.foobar"
	alertLabel := acctest.RandomWithPrefix("tf-test")

	// TODO: revert to use alert channels from API once it's available
	channelID := 10000

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilter(t, alertLabel, channelID),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resName, "alert_definitions.#", 0),
					acceptance.CheckResourceAttrContains(resName, "alert_definitions.0.label", alertLabel),
				),
			},
		},
	})
}
