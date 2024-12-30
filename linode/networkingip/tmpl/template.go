package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Region   string
	Reserved bool
}

func DataBasic(t testing.TB, label, region string) string {
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

func NetworkingIPReserveTest(t *testing.T, label string, region string, reserved bool) string {
	return acceptance.ExecuteTemplate(t, "networking_ip_reserve_test", TemplateData{
		Label:    label,
		Region:   region,
		Reserved: reserved,
	})
}
