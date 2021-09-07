package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Image    string
	ID       string
	FilePath string
}

func Basic(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"image_basic", TemplateData{Image: image})
}

func Updates(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"image_updates", TemplateData{Image: image})
}

func Upload(t *testing.T, image, upload string) string {
	return acceptance.ExecuteTemplate(t,
		"image_upload", TemplateData{
			Image:    image,
			FilePath: upload,
		})
}

func DataBasic(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"image_data_basic", TemplateData{ID: id})
}
