package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label string
	VPCId int
	IPv4  string
}

func Basic(t *testing.T, vpcId int, label, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_basic", TemplateData{
			Label: label,
			VPCId: vpcId,
			IPv4:  ipv4,
		})
}

func Updates(t *testing.T, vpcId int, label, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_updates", TemplateData{
			Label: label,
			VPCId: vpcId,
			IPv4:  ipv4,
		})
}

func DataBasic(t *testing.T, vpcId int, label, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_data_basic", TemplateData{
			Label: label,
			VPCId: vpcId,
			IPv4:  ipv4,
		})
}
