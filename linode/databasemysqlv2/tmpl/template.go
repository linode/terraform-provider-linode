package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
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
	Updates     TemplateDataUpdates
}

type TemplateDataEngineConfig struct {
	Label    string
	Region   string
	EngineID string
	Type     string

	EngineConfigBinlogRetentionPeriod             int
	EngineConfigMySQLConnectTimeout               int
	EngineConfigMySQLDefaultTimeZone              string
	EngineConfigMySQLGroupConcatMaxLen            float64
	EngineConfigMySQLInformationSchemaStatsExpiry int
	EngineConfigMySQLInnoDBChangeBufferMaxSize    int
	EngineConfigMySQLInnoDBFlushNeighbors         int
	EngineConfigMySQLInnoDBFTMinTokenSize         int
	EngineConfigMySQLInnoDBFTServerStopwordTable  string
	EngineConfigMySQLInnoDBLockWaitTimeout        int
	EngineConfigMySQLInnoDBLogBufferSize          int
	EngineConfigMySQLInnoDBOnlineAlterLogMaxSize  int
	EngineConfigMySQLInnoDBReadIOThreads          int
	EngineConfigMySQLInnoDBRollbackOnTimeout      bool
	EngineConfigMySQLInnoDBThreadConcurrency      int
	EngineConfigMySQLInnoDBWriteIOThreads         int
	EngineConfigMySQLInteractiveTimeout           int
	EngineConfigMySQLInternalTmpMemStorageEngine  string
	EngineConfigMySQLMaxAllowedPacket             int
	EngineConfigMySQLMaxHeapTableSize             int
	EngineConfigMySQLNetBufferLength              int
	EngineConfigMySQLNetReadTimeout               int
	EngineConfigMySQLNetWriteTimeout              int
	EngineConfigMySQLSortBufferSize               int
	EngineConfigMySQLSQLMode                      string
	EngineConfigMySQLSQLRequirePrimaryKey         bool
	EngineConfigMySQLTmpTableSize                 int
	EngineConfigMySQLWaitTimeout                  int
}

func Basic(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_basic",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func Complex(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_complex",
		data,
	)
}

func EngineConfig(
	t testing.TB,
	data TemplateDataEngineConfig,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_engine_config",
		data,
	)
}

func EngineConfigNullableField(
	t testing.TB,
	data TemplateDataEngineConfig,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_engine_config_nullable_field",
		data,
	)
}

func Fork(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_fork",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func Suspension(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_suspension",
		data,
	)
}

func VPC0(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_vpc_0",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func VPC1(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_vpc_1",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}

func Data(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_data",
		data,
	)
}

func DataEngineConfig(
	t testing.TB,
	data TemplateDataEngineConfig,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_data_engine_config",
		data,
	)
}

func DataVPC(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_mysql_v2_data_vpc",
		TemplateData{
			Label:    label,
			Region:   region,
			EngineID: engine,
			Type:     nodeType,
		},
	)
}
