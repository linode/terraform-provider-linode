//go:build integration || monitoralertdefinition

package monitoralertdefinition_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinition/tmpl"
)

var channelID int

func init() {
	resource.AddTestSweepers("linode_monitor_alert_definition", &resource.Sweeper{
		Name: "linode_monitor_alert_definition",
		F:    sweep,
	})
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
	//if len(channels) < 1 {
	//	fmt.Errorf("at least one alert channel is required for alert definition tests")
	//}
	//
	//channelID = channels[0].ID

	channelID = 10000
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	alertDefinitions, err := client.ListAllMonitorAlertDefinitions(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("error getting alert definitions: %s", err)
	}
	for _, alert := range alertDefinitions {
		if !acceptance.ShouldSweep(prefix, alert.Label) {
			continue
		}
		err := client.DeleteMonitorAlertDefinition(context.Background(), alert.ServiceType, alert.ID)
		if err != nil {
			return fmt.Errorf("error destroying %v during sweep: %s", alert.Label, err)
		}
	}

	return nil
}

func TestAccResourceAlertDefinition_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_monitor_alert_definition.test"
	alertLabel := acctest.RandomWithPrefix("tf-test")
	alertChannels := channelID
	aggregateFunction := "avg"
	aggregateFunctionUpdate := "sum"
	triggerOccurrences := 1
	triggerOccurrencesUpdate := 3

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkAlertDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, alertLabel, aggregateFunction, alertChannels, triggerOccurrences),
				Check:  checkAlertDefinitionExists,
				ConfigStateChecks: []statecheck.StateCheck{
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
						knownvalue.StringExact(aggregateFunction),
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
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("trigger_conditions").AtMapKey("trigger_occurrences"),
						knownvalue.Int64Exact(int64(triggerOccurrences)),
					),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       resourceImportStateID,
				ImportStateVerifyIgnore: []string{"wait_for"},
			},
			{
				Config: tmpl.Updates(t, alertLabel, aggregateFunctionUpdate, alertChannels, triggerOccurrencesUpdate),
				Check:  checkAlertDefinitionExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("label"), knownvalue.StringExact(fmt.Sprintf("%s-updated", alertLabel))),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("trigger_conditions").AtMapKey("trigger_occurrences"),
						knownvalue.Int64Exact(int64(triggerOccurrencesUpdate)),
					),
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("rule_criteria").AtMapKey("rules").AtSliceIndex(0).AtMapKey("aggregate_function"),
						knownvalue.StringExact(aggregateFunctionUpdate),
					),
				},
			},
		},
	})
}

func checkAlertDefinitionExists(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_monitor_alert_definition" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}

		serviceType := rs.Primary.Attributes["service_type"]

		_, err = client.GetMonitorAlertDefinition(context.Background(), serviceType, id)
		if err != nil {
			return fmt.Errorf("error retrieving state of Alert Definition %d: %s", id, err)
		}
	}

	return nil
}

func checkAlertDefinitionDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_monitor_alert_definition" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("would have considered %v as %d", rs.Primary.ID, id)
		}

		serviceType := rs.Primary.Attributes["service_type"]

		_, err = client.GetMonitorAlertDefinition(context.Background(), serviceType, id)

		if err == nil {
			return fmt.Errorf("Alert Definition with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting Alert Definition with id %d", id)
		}
	}

	return nil
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_monitor_alert_definition" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}

		serviceType := rs.Primary.Attributes["service_type"]

		return fmt.Sprintf("%d,%s", id, serviceType), nil
	}

	return "", fmt.Errorf("Error finding linode_monitor_alert_definition")
}
