package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label        string
	SpaceID      int
	FilesystemID int
}

// Basic renders a quota with only a capacity limit.
func Basic(t testing.TB, label string, spaceID, filesystemID int) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_basic", TemplateData{
			Label:        label,
			SpaceID:      spaceID,
			FilesystemID: filesystemID,
		})
}

// Updates renders a quota with default user and group capacity limits.
func Updates(t testing.TB, label string, spaceID, filesystemID int) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_updates", TemplateData{
			Label:        label,
			SpaceID:      spaceID,
			FilesystemID: filesystemID,
		})
}

func DataBasic(t testing.TB, label string, spaceID, filesystemID int) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_data_basic", TemplateData{
			Label:        label,
			SpaceID:      spaceID,
			FilesystemID: filesystemID,
		})
}
