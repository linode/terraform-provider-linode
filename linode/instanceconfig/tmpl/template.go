package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Booted bool
	Swap   bool
	Region string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Complex(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex", TemplateData{
			Label:  label,
			Region: region,
		})
}

func ComplexUpdates(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex_updates", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Booted(t *testing.T, label, region string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted", TemplateData{
			Label:  label,
			Booted: booted,
			Region: region,
		})
}

func BootedSwap(t *testing.T, label, region string, swap bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted_swap", TemplateData{
			Label:  label,
			Swap:   swap,
			Region: region,
		})
}

func Provisioner(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_provisioner", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DeviceBlock(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(
		t,
		"instance_config_device_block", TemplateData{
			Label:  label,
			Region: region,
		},
	)
}

func DeviceNamedBlock(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(
		t,
		"instance_config_device_named_block", TemplateData{
			Label:  label,
			Region: region,
		},
	)
}
