package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func RouteTarget(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_route_target", TemplateData{
			Label:  label,
			Region: region,
		})
}

func NoID(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_no_id", nil)
}

func ReassignmentStep1(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_reassign_step1", TemplateData{
			Label:  label,
			Region: region,
		})
}

func ReassignmentStep2(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_reassign_step2", TemplateData{
			Label:  label,
			Region: region,
		})
}

func RaceCondition(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_race_condition", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}
