package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_data_basic", TemplateData{Label: label, Region: region})
}

func NetworkingIPReservedAssigned(t *testing.T, label string, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_reserved_assigned",
		TemplateData{
			Label:  label,
			Region: region,
		})
}

func NetworkingIPReserved(t *testing.T, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_reserved",
		TemplateData{
			Region: region,
		})
}

func NetworkingIPEphemeral(t *testing.T, label string, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_ephemeral",
		TemplateData{
			Label:  label,
			Region: region,
		})
}

func NetworkingIPEphemeralAssigned(t *testing.T, label string, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_ephemeral_assigned",
		TemplateData{
			Label:  label,
			Region: region,
		})
}
