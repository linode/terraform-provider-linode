package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Country      string
	Status       string
	Capabilities string
}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"regions_data_basic", nil)
}

func DataFilterCountry(t *testing.T, country, status string, capabilities string) string {
	return acceptance.ExecuteTemplate(t,
		"regions_data_filter_by_country", TemplateData{
			Country:      country,
			Status:       status,
			Capabilities: capabilities,
		})
}

func DataFilterStatus(t *testing.T, country, status string, capabilities string) string {
	return acceptance.ExecuteTemplate(t,
		"regions_data_filter_by_status", TemplateData{
			Country:      country,
			Status:       status,
			Capabilities: capabilities,
		})
}

func DataFilterCapabilities(t *testing.T, country, status string, capabilities string) string {
	return acceptance.ExecuteTemplate(t,
		"regions_data_filter_by_capabilities", TemplateData{
			Country:      country,
			Status:       status,
			Capabilities: capabilities,
		})
}
