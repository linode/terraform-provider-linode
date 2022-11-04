package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Image    string
	ID       string
	FilePath string
	Region   string
}

func Basic(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"image_basic", TemplateData{
			Image:  image,
			Region: region,
		})
}

func Updates(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"image_updates", TemplateData{
			Image:  image,
			Region: region,
		})
}

func Upload(t *testing.T, image, upload, region string) string {
	return acceptance.ExecuteTemplate(t,
		"image_upload", TemplateData{
			Image:    image,
			FilePath: upload,
			Region:   region,
		})
}

func DataBasic(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"image_data_basic", TemplateData{ID: id})
}
