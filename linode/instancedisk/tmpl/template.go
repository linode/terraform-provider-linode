package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Size     int
	PubKey   string
	Region   string
	RootPass string
}

func Basic(t *testing.T, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label:  label,
			Size:   size,
			Region: region,
		})
}

func Complex(t *testing.T, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label:  label,
			Size:   size,
			PubKey: acceptance.PublicKeyMaterial,
			Region: region,
		})
}

func BootedResize(t *testing.T, label, region string, size int, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_booted_resize", TemplateData{
			Label:    label,
			Size:     size,
			Region:   region,
			RootPass: rootPass,
		})
}
