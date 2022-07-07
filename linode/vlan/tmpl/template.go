package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	InstLabel string
	VLANLabel string
}

func DataBasic(t *testing.T, instLabel, vlanLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_basic", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
		})
}

func DataRegex(t *testing.T, instLabel, vlanLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_regex", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
		})
}

func DataCheckDuplicate(t *testing.T, instLabel, vlanLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"vlan_data_check_duplicate", TemplateData{
			InstLabel: instLabel,
			VLANLabel: vlanLabel,
		})
}
