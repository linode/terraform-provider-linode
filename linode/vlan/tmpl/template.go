package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	InstLabel string
	VLANLabel string
	Region    string
	Label     string
	RootPass  string
}

func DataBasic(t testing.TB, instLabel, region, vlanLabel, label, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_basic", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
			RootPass:  rootPass,
		})
}

func DataRegex(t testing.TB, instLabel, region, vlanLabel, label, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_regex", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
			RootPass:  rootPass,
		})
}

func DataCheckDuplicate(t testing.TB, instLabel, region, vlanLabel, label, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_check_duplicate", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
			Region:    region,
			Label:     label,
			RootPass:  rootPass,
		})
}
