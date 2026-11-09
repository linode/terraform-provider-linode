package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	SpaceID          int
	FilesystemID     int
	Path             string
	MaxCapacityBytes int64
	MaxFileCount     int64
}

func Basic(t testing.TB, spaceID, filesystemID int, path string, maxCapacityBytes, maxFileCount int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_basic", TemplateData{
			SpaceID:          spaceID,
			FilesystemID:     filesystemID,
			Path:             path,
			MaxCapacityBytes: maxCapacityBytes,
			MaxFileCount:     maxFileCount,
		})
}

func Updates(t testing.TB, spaceID, filesystemID int, path string, maxCapacityBytes, maxFileCount int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_updates", TemplateData{
			SpaceID:          spaceID,
			FilesystemID:     filesystemID,
			Path:             path,
			MaxCapacityBytes: maxCapacityBytes,
			MaxFileCount:     maxFileCount,
		})
}

func DataBasic(t testing.TB, spaceID, filesystemID int, path string, maxCapacityBytes, maxFileCount int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_quota_data_basic", TemplateData{
			SpaceID:          spaceID,
			FilesystemID:     filesystemID,
			Path:             path,
			MaxCapacityBytes: maxCapacityBytes,
			MaxFileCount:     maxFileCount,
		})
}
