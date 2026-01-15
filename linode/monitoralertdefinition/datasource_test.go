//go:build integration || monitoralertdefinition

package monitoralertdefinition_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinition/tmpl"
)

func TestAccDataSourceAlertDefinition_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_monitor_alert_definition.foobar"
	alertLabel := acctest.RandomWithPrefix("tf-test")
	alertChannels := channelID

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkAlertDefinitionDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, alertLabel, "avg", alertChannels, 1),
				Check:  checkAlertDefinitionExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("service_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("severity"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("channel_ids"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("entity_ids"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("has_more_resources"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("created_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("updated_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("class"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_channels").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_channels").AtSliceIndex(0).AtMapKey("label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_channels").AtSliceIndex(0).AtMapKey("type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_channels").AtSliceIndex(0).AtMapKey("url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("aggregate_function"),
						knownvalue.StringExact("avg"),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("metric"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("operator"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("threshold"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("dimension_filters").
							AtSliceIndex(0).
							AtMapKey("dimension_label"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("dimension_filters").AtSliceIndex(0).AtMapKey("operator"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("dimension_filters").AtSliceIndex(0).AtMapKey("value"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("trigger_conditions").AtMapKey("criteria_condition"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("trigger_conditions").AtMapKey("evaluation_period_seconds"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("trigger_conditions").AtMapKey("polling_interval_seconds"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("trigger_conditions").AtMapKey("trigger_occurrences"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
