package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	FirewallLabel   string
	InstanceLabel   string
	InstanceRegion  string
	ImageLabel1     string
	ImageLabel2     string
	ShareGroupLabel string
	TokenLabel      string
	MemberLabel     string
}

func DataBasic(t testing.TB, fwLabel, instanceLabel, instanceRegion, imageLabel1, imageLabel2, shareGroupLabel, tokenLabel, memberLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"consumer_image_share_group_image_shares_data_basic", TemplateData{
			FirewallLabel:   fwLabel,
			InstanceLabel:   instanceLabel,
			InstanceRegion:  instanceRegion,
			ImageLabel1:     imageLabel1,
			ImageLabel2:     imageLabel2,
			ShareGroupLabel: shareGroupLabel,
			TokenLabel:      tokenLabel,
			MemberLabel:     memberLabel,
		})
}
