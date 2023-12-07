package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Booted   bool
	Swap     bool
	Region   string
	RootPass string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Complex(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}

func ComplexUpdates(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex_updates", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}

func Booted(t *testing.T, label, region string, booted bool, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted", TemplateData{
			Label:    label,
			Booted:   booted,
			Region:   region,
			RootPass: rootPass,
		})
}

func BootedSwap(t *testing.T, label, region string, swap bool, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted_swap", TemplateData{
			Label:    label,
			Swap:     swap,
			Region:   region,
			RootPass: rootPass,
		})
}

func Provisioner(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_provisioner", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}

func DeviceBlock(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(
		t,
		"instance_config_device_block", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		},
	)
}

func DeviceNamedBlock(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(
		t,
		"instance_config_device_named_block", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		},
	)
}

func VPCInterface(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(
		t,
		"vpc_interface", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		},
	)
}

func VPCInterfaceUpdates(t *testing.T, label, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(
		t,
		"vpc_interface_update", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		},
	)
}
