package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	objkey "github.com/linode/terraform-provider-linode/v2/linode/objkey/tmpl"
)

type TemplateData struct {
	Key objkey.TemplateData

	Label       string
	ACL         string
	CORSEnabled bool
	Versioning  bool

	Cert    string
	PrivKey string
	Cluster string
	Region  string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_basic", TemplateData{Label: label, Region: region})
}

func BasicLegacy(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_basic", TemplateData{Label: label, Cluster: cluster})
}

func Updates(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_updates", TemplateData{Label: label, Region: region})
}

func Access(t *testing.T, label, cluster, acl string, cors bool) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_access", TemplateData{
			Label:       label,
			ACL:         acl,
			CORSEnabled: cors,
			Cluster:     cluster,
		})
}

func Cert(t *testing.T, label, region, cert, privKey string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_cert", TemplateData{
			Label:   label,
			Cert:    cert,
			PrivKey: privKey,
			Region:  region,
		})
}

func Versioning(t *testing.T, label, cluster, keyName string, versioning bool) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_versioning", TemplateData{
			Key:        objkey.TemplateData{Label: keyName},
			Label:      label,
			Versioning: versioning,
			Cluster:    cluster,
		})
}

func LifeCycle(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func LifeCycleNoID(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle_no_id", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func LifeCycleUpdates(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle_updates", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func LifeCycleRemoved(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle_removed", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func TempKeys(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_temp_keys", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func ForceDelete(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_force_delete", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func ForceDelete_Empty(t *testing.T) string {
	return acceptance.ExecuteTemplate(t, "object_bucket_force_delete_empty", nil)
}

func ClusterDataBasic(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_cluster_data_basic", TemplateData{
			Label:   label,
			Cluster: cluster,
		})
}

func CredsConfiged(t *testing.T, label, cluster, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_creds_configed", TemplateData{
			Key:     objkey.TemplateData{Label: keyName},
			Label:   label,
			Cluster: cluster,
		})
}

func DataBasic(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_data_basic", TemplateData{
			Label:   label,
			Cluster: cluster,
		})
}
