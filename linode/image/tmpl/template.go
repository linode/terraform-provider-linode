package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Image         string
	ID            string
	FilePath      string
	Region        string
	Label         string
	Tag           string
	ReplicaRegion string
}

func Basic(t *testing.T, image, region, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"image_basic", TemplateData{
			Image:  image,
			Region: region,
			Label:  label,
			Tag:    tag,
		})
}

func Updates(t *testing.T, image, region, label, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"image_updates", TemplateData{
			Image:  image,
			Region: region,
			Label:  label,
			Tag:    tag,
		})
}

func Upload(t *testing.T, image, upload, region, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"image_upload", TemplateData{
			Image:    image,
			FilePath: upload,
			Region:   region,
			Tag:      tag,
		})
}

func Replicate(t *testing.T, image, upload, region, replicaRegion string) string {
	return acceptance.ExecuteTemplate(t,
		"image_replicate", TemplateData{
			Image:         image,
			Region:        region,
			FilePath:      upload,
			ReplicaRegion: replicaRegion,
		})
}

func NoReplicaRegions(t *testing.T, image, upload, region string) string {
	return acceptance.ExecuteTemplate(t,
		"image_no_replica_regions", TemplateData{
			Image:    image,
			Region:   region,
			FilePath: upload,
		})
}

func DataBasic(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"image_data_basic", TemplateData{ID: id})
}

func DataReplicate(t *testing.T, image, upload, region, replicaRegion string) string {
	return acceptance.ExecuteTemplate(t,
		"image_data_replicate", TemplateData{
			Image:         image,
			Region:        region,
			FilePath:      upload,
			ReplicaRegion: replicaRegion,
		})
}
