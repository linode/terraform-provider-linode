package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label             string
	SpaceID           int
	FilesystemID      int
	MaxFileCount      int64
	UserMaxFileCount  int64
	GroupMaxFileCount int64
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

// Updates renders a quota with user and group limits; the default user and
// group file count limits are a quarter and a half of maxFileCount.
func Updates(t testing.TB, label string, spaceID, filesystemID int, maxFileCount int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_updates", TemplateData{
			Label:             label,
			SpaceID:           spaceID,
			FilesystemID:      filesystemID,
			MaxFileCount:      maxFileCount,
			UserMaxFileCount:  maxFileCount / 4,
			GroupMaxFileCount: maxFileCount / 2,
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
