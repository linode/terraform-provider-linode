package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label                         string
	Region                        string
	EngineID                      string
	Type                          string
	AllowedIP                     string
	ClusterSize                   int
	MaintenanceHour               int
	Suspended                     bool
	BackupHour                    int
	BackupMinute                  int
	Timeout                       int
	MaxmemoryPolicy               string
	FrequentSnapshots             bool
	ACLChannelsDefault            string
	ActiveExpireEffort            int
	ActiveDefrag                  bool
	LFUDecayTime                  int
	LFULogFactor                  int
	NumberOfDatabases             int
	Persistence                   string
	PubsubClientOutputBufferLimit int
}

func Basic(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_basic", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func DataBasic(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_basic", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func Complex(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_complex", data)
}

func Suspension(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_suspension", data)
}

func EngineConfig(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_engine_config", data)
}

func Resize(t testing.TB, label, region, engineID, nodeType string, clusterSize int) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_resize", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType, ClusterSize: clusterSize,
	})
}

func Fork(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_fork", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func VPC0(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_vpc_0", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func VPC1(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_vpc_1", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func VPCDetach(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_vpc_detach", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}

func DataEngineConfig(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_engine_config", data)
}

func DataVPC(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_vpc", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}
