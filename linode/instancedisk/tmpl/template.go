package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Size   int
	PubKey string
}

func Basic(t *testing.T, label string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label: label,
			Size:  size,
		})
}

func Complex(t *testing.T, label string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label:  label,
			Size:   size,
			PubKey: acceptance.PublicKeyMaterial,
		})
}

func BootedResize(t *testing.T, label string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_booted_resize", TemplateData{
			Label: label,
			Size:  size,
		})
}
