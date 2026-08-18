package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label               string
	Region              string
	RootPass            string
	AssignedLinodeIndex int
	Reserved            bool
}

func DataBasic(t testing.TB, label, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ip_data_basic", TemplateData{Label: label, Region: region, RootPass: rootPass})
}

func NetworkingIPReservedUnassigned(t *testing.T, region string, reserved bool) string {
	return acceptance.ExecuteTemplate(t, "networking_ip_reserved_unassigned", TemplateData{
		Region:   region,
		Reserved: reserved,
	})
}

func NetworkingIPReservedAssigned(
	t *testing.T,
	label string,
	region string,
	assignedLinodeIndex int,
	reserved bool,
) string {
	return acceptance.ExecuteTemplate(t, "networking_ip_reserved_assigned", TemplateData{
		Label:               label,
		Region:              region,
		Reserved:            reserved,
		AssignedLinodeIndex: assignedLinodeIndex,
	})
}
