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
}

func Basic(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_basic", TemplateData{Label: label, Cluster: cluster})
}

func Updates(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_updates", TemplateData{Label: label, Cluster: cluster})
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

func Cert(t *testing.T, label, cluster, cert, privKey string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_cert", TemplateData{
			Label:   label,
			Cert:    cert,
			PrivKey: privKey,
			Cluster: cluster,
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

func ClusterDataBasic(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_cluster_data_basic", TemplateData{
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
