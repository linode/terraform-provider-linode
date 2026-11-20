package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateDataUpdates struct {
	HourOfDay, DayOfWeek, Duration int
	Frequency                      string
}

type TemplateData struct {
	Label       string
	Region      string
	EngineID    string
	Type        string
	AllowedIP   string
	ClusterSize int
	Suspended   bool

	Updates TemplateDataUpdates
}

type TemplateDataEngineConfig struct {
	Label    string
	Region   string
	EngineID string
	Type     string

	EngineConfigBackupHour                          int
	EngineConfigBackupMinute                        int
	EngineConfigFrequentSnapshots                   bool
	EngineConfigValkeyACLChannelsDefault            string
	EngineConfigValkeyActiveExpireEffort            int
	EngineConfigValkeyActiveDefrag                  bool
	EngineConfigValkeyLFUDecayTime                  int
	EngineConfigValkeyLFULogFactor                  int
	EngineConfigValkeyMaxmemoryPolicy               string
	EngineConfigValkeyNumberOfDatabases             int
	EngineConfigValkeyPersistence                   string
	EngineConfigValkeyPubsubClientOutputBufferLimit int
	EngineConfigValkeyTimeout                       int
}

func Basic(t testing.TB, label, region, engineID string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_basic", TemplateData{
		Label: label, Region: region, EngineID: engineID,
	})
}

func Data(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data", data)
}

func DataListing(t testing.TB, label, region, engineID string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_listing", TemplateData{
		Label: label, Region: region, EngineID: engineID,
	})
}

func Complex(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_complex", data)
}

func Suspension(t testing.TB, data TemplateData) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_suspension", data)
}

func EngineConfig(t testing.TB, data TemplateDataEngineConfig) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_engine_config", data)
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

func DataEngineConfig(t testing.TB, data TemplateDataEngineConfig) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_engine_config", data)
}

func DataVPC(t testing.TB, label, region, engineID, nodeType string) string {
	return acceptance.ExecuteTemplate(t, "database_valkey_v2_data_vpc", TemplateData{
		Label: label, Region: region, EngineID: engineID, Type: nodeType,
	})
}
