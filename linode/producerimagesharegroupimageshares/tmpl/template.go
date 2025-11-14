package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label           string
	InstanceLabel   string
	InstanceRegion  string
	ImageLabel1     string
	ImageLabel2     string
	ShareGroupLabel string
}

func DataBasic(t testing.TB, label, instanceLabel, instanceRegion, imageLabel1, imageLabel2, shareGroupLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_group_image_shares_data_basic", TemplateData{
			Label:           label,
			InstanceLabel:   instanceLabel,
			InstanceRegion:  instanceRegion,
			ImageLabel1:     imageLabel1,
			ImageLabel2:     imageLabel2,
			ShareGroupLabel: shareGroupLabel,
		})
}
