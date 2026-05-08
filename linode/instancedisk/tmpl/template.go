package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Size     int
	PubKey   string
	Region   string
	Image    string
	RootPass string
}

func Basic(t testing.TB, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label:  label,
			Size:   size,
			Region: region,
		})
}

func Complex(t testing.TB, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_basic", TemplateData{
			Label:  label,
			Size:   size,
			PubKey: acceptance.PublicKeyMaterial,
			Region: region,
		})
}

func BootedResize(t testing.TB, label, region string, size int, rootPass string) string {
	return BootedResizeWithImage(t, label, region, size, "linode/debian13", rootPass)
}

func BootedResizeWithImage(t testing.TB, label, region string, size int, image, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_booted_resize", TemplateData{
			Label:    label,
			Size:     size,
			Region:   region,
			Image:    image,
			RootPass: rootPass,
		})
}

func ImageAuthKeysOnly(t testing.TB, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_image_auth_keys_only", TemplateData{
			Label:  label,
			Size:   size,
			PubKey: acceptance.PublicKeyMaterial,
			Region: region,
		})
}

func ImageNoAuth(t testing.TB, label, region string, size int) string {
	return acceptance.ExecuteTemplate(t,
		"instance_disk_image_no_auth", TemplateData{
			Label:  label,
			Size:   size,
			Region: region,
		})
}
