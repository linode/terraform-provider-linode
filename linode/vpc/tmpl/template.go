package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
<<<<<<< HEAD
	Label   string
	Region  string
	VPCType string
=======
	Label      string
	Region     string
	IPv4Range  string
	IPv4Range2 string
	VPCType    string
>>>>>>> dev
}

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Updates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_updates", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DualStack(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_dual_stack", TemplateData{
			Label:  label,
			Region: region,
		})
}

func IPv4(t testing.TB, label, region, ipv4Range string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_ipv4", TemplateData{
			Label:     label,
			Region:    region,
			IPv4Range: ipv4Range,
		})
}

func IPv4Update(t testing.TB, label, region, ipv4Range, ipv4Range2 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_ipv4_update", TemplateData{
			Label:      label,
			Region:     region,
			IPv4Range:  ipv4Range,
			IPv4Range2: ipv4Range2,
		})
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataDualStack(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_data_dual_stack", TemplateData{
			Label:  label,
			Region: region,
		})
}

<<<<<<< HEAD
=======
func DataIPv4(t testing.TB, label, region, ipv4Range string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_data_ipv4", TemplateData{
			Label:     label,
			Region:    region,
			IPv4Range: ipv4Range,
		})
}

>>>>>>> dev
func VPCType(t testing.TB, label, region, vpcType string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_vpc_type", TemplateData{
			Label:   label,
			Region:  region,
			VPCType: vpcType,
		})
}
