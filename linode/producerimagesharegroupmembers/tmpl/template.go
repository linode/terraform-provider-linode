package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	ShareGroupLabel string
	TokenLabel      string
	MemberLabel     string
}

func DataBasic(t testing.TB, shareGroupLabel, tokenLabel, memberLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_group_members_data_basic", TemplateData{
			ShareGroupLabel: shareGroupLabel,
			TokenLabel:      tokenLabel,
			MemberLabel:     memberLabel,
		})
}
