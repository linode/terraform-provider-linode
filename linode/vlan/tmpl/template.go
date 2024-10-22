package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	InstLabel string
	VLANLabel string
	Region    string
	Label     string
}

func DataBasic(t testing.TB, instLabel, region, vlanLabel, label string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_basic", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
		})
}

func DataRegex(t testing.TB, instLabel, region, vlanLabel, label string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_regex", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
		})
}

func DataCheckDuplicate(t testing.TB, instLabel, region, vlanLabel, label string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_check_duplicate", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
		})
}
