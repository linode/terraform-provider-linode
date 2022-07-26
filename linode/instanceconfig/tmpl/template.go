package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Booted bool
	Swap   bool
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_basic", TemplateData{
			Label: label,
		})
}

func Complex(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex", TemplateData{
			Label: label,
		})
}

func ComplexUpdates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex_updates", TemplateData{
			Label: label,
		})
}

func Booted(t *testing.T, label string, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted", TemplateData{
			Label:  label,
			Booted: booted,
		})
}

func BootedSwap(t *testing.T, label string, swap bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_booted_swap", TemplateData{
			Label: label,
			Swap:  swap,
		})
}

func DiskSwap(t *testing.T, label string, swap, booted bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_disk_swap", TemplateData{
			Label:  label,
			Swap:   swap,
			Booted: booted,
		})
}
