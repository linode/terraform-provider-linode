package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_basic", TemplateData{
			Label: label,
		})
}

func RouteTarget(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_route_target", TemplateData{
			Label: label,
		})
}

func NoID(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_no_id", nil)
}

func ReassignmentStep1(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_reassign_step1", TemplateData{
			Label: label,
		})
}

func ReassignmentStep2(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_reassign_step2", TemplateData{
			Label: label,
		})
}

func RaceCondition(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_race_condition", TemplateData{
			Label: label,
		})
}

func DataBasic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_data_basic", TemplateData{
			Label: label,
		})
}
