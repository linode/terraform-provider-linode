package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label       string
	Description string
}

type UpdateTemplateData struct {
	Label                      string
	Region                     string
	ImageLabel1                string
	ImageLabel2                string
	ImageShareGroupLabel       string
	ImageShareGroupDescription string
	Images                     []ShareGroupImageTemplate
}

type ShareGroupImageTemplate struct {
	ID          string
	Label       string
	Description string
}

func DataBasic(t testing.TB, label, description string) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_group_data_basic", TemplateData{
			Label:       label,
			Description: description,
		})
}

func Basic(t testing.TB, label, description string) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_group_basic", TemplateData{
			Label:       label,
			Description: description,
		})
}

func Updates(t testing.TB, label, region, image_label_1, image_label_2, isg_label, isg_description string, images []ShareGroupImageTemplate) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_group_updates", UpdateTemplateData{
			Label:                      label,
			Region:                     region,
			ImageLabel1:                image_label_1,
			ImageLabel2:                image_label_2,
			ImageShareGroupLabel:       isg_label,
			ImageShareGroupDescription: isg_description,
			Images:                     images,
		})
}
