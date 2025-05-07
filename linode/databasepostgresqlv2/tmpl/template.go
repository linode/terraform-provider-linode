package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
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

	EngineConfigPGAutovacuumAnalyzeScaleFactor         float64
	EngineConfigPGAutovacuumAnalyzeThreshold           int
	EngineConfigPGAutovacuumMaxWorkers                 int
	EngineConfigPGAutovacuumNaptime                    int
	EngineConfigPGAutovacuumVacuumCostDelay            int
	EngineConfigPGAutovacuumVacuumCostLimit            int
	EngineConfigPGAutovacuumVacuumScaleFactor          float64
	EngineConfigPGAutovacuumVacuumThreshold            int
	EngineConfigPGBGWriterDelay                        int
	EngineConfigPGBGWriterFlushAfter                   int
	EngineConfigPGBGWriterLRUMaxpages                  int
	EngineConfigPGBGWriterLRUMultiplier                float64
	EngineConfigPGDeadlockTimeout                      int
	EngineConfigPGDefaultToastCompression              string
	EngineConfigPGIdleInTransactionSessionTimeout      int
	EngineConfigPGJIT                                  bool
	EngineConfigPGMaxFilesPerProcess                   int
	EngineConfigPGMaxLocksPerTransaction               int
	EngineConfigPGMaxLogicalReplicationWorkers         int
	EngineConfigPGMaxParallelWorkers                   int
	EngineConfigPGMaxParallelWorkersPerGather          int
	EngineConfigPGMaxPredLocksPerTransaction           int
	EngineConfigPGMaxReplicationSlots                  int
	EngineConfigPGMaxSlotWALKeepSize                   int
	EngineConfigPGMaxStackDepth                        int
	EngineConfigPGMaxStandbyArchiveDelay               int
	EngineConfigPGMaxStandbyStreamingDelay             int
	EngineConfigPGMaxWALSenders                        int
	EngineConfigPGMaxWorkerProcesses                   int
	EngineConfigPGPasswordEncryption                   string
	EngineConfigPGPGPartmanBGWInterval                 int
	EngineConfigPGPGPartmanBGWRole                     string
	EngineConfigPGPGStatMonitorPGSMEnableQueryPlan     bool
	EngineConfigPGPGStatMonitorPGSMMaxBuckets          int
	EngineConfigPGPGStatStatementsTrack                string
	EngineConfigPGTempFileLimit                        int
	EngineConfigPGTimezone                             string
	EngineConfigPGTrackActivityQuerySize               int
	EngineConfigPGTrackCommitTimestamp                 string
	EngineConfigPGTrackFunctions                       string
	EngineConfigPGTrackIOTiming                        string
	EngineConfigPGWALSenderTimeout                     int
	EngineConfigPGWALWriterDelay                       int
	EngineConfigPGStatMonitorEnable                    bool
	EngineConfigPGLookoutMaxFailoverReplicationTimeLag int
	EngineConfigSharedBuffersPercentage                float64
	EngineConfigWorkMem                                int
}

func Basic(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_basic",
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
		"database_postgresql_v2_complex",
		data,
	)
}

func Fork(t testing.TB, label, region, engine, nodeType string) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_fork",
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
		"database_postgresql_v2_suspension",
		data,
	)
}

func Data(
	t testing.TB,
	data TemplateData,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_data",
		data,
	)
}

func EngineConfig(
	t testing.TB,
	data TemplateDataEngineConfig,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_engine_config",
		data,
	)
}

func DataEngineConfig(
	t testing.TB,
	data TemplateDataEngineConfig,
) string {
	return acceptance.ExecuteTemplate(
		t,
		"database_postgresql_v2_data_engine_config",
		data,
	)
}
