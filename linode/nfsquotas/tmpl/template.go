package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	SpaceID      int
	FilesystemID int
}

func DataBasic(t testing.TB, spaceID, filesystemID int) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quotas_data_basic", TemplateData{
			SpaceID:      spaceID,
			FilesystemID: filesystemID,
		})
}
