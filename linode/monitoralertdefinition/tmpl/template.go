package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label              string
	AlertChannel       int
	AggregateFunction  string
	TriggerOccurrences int
	GroupBy            []string
}

func Basic(t testing.TB, label, aggregateFunction string, alertChannel, triggerOccurrences int, groupBy ...string) string {
	return acceptance.ExecuteTemplate(t,
		"alert_definition_basic", TemplateData{
			Label:              label,
			AlertChannel:       alertChannel,
			AggregateFunction:  aggregateFunction,
			TriggerOccurrences: triggerOccurrences,
			GroupBy:            groupBy,
		})
}

func DataBasic(t testing.TB, label, aggregateFunction string, alertChannel, triggerOccurrences int, groupBy ...string) string {
	return acceptance.ExecuteTemplate(t,
		"alert_definition_data_basic", TemplateData{
			Label:              label,
			AlertChannel:       alertChannel,
			AggregateFunction:  aggregateFunction,
			TriggerOccurrences: triggerOccurrences,
			GroupBy:            groupBy,
		})
}

func Updates(t testing.TB, label, aggregateFunction string, alertChannel, triggerOccurrences int, groupBy ...string) string {
	return acceptance.ExecuteTemplate(t,
		"alert_definition_update", TemplateData{
			Label:              label,
			AlertChannel:       alertChannel,
			AggregateFunction:  aggregateFunction,
			TriggerOccurrences: triggerOccurrences,
			GroupBy:            groupBy,
		})
}
