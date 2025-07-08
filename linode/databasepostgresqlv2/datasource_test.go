//go:build integration || databasepostgresqlv2

package databasepostgresqlv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/databasepostgresqlv2/tmpl"
)

func TestAccDataSource_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_database_postgresql_v2.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Data(t, tmpl.TemplateData{
					Label:       label,
					Region:      testRegion,
					EngineID:    testEngine,
					Type:        "g6-nanode-1",
					AllowedIP:   "10.0.0.3/32",
					ClusterSize: 1,
					Updates: tmpl.TemplateDataUpdates{
						HourOfDay: 3,
						DayOfWeek: 2,
						Duration:  4,
						Frequency: "weekly",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),

					resource.TestCheckResourceAttrSet(dataSourceName, "ca_cert"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttr(dataSourceName, "encrypted", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(dataSourceName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(dataSourceName, "fork_source"),
					resource.TestCheckResourceAttrSet(dataSourceName, "host_primary"),
					resource.TestCheckResourceAttr(dataSourceName, "label", label),
					resource.TestCheckResourceAttrSet(dataSourceName, "members.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "root_password"),
					resource.TestCheckResourceAttrSet(dataSourceName, "root_username"),
					resource.TestCheckResourceAttr(dataSourceName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "port"),
					resource.TestCheckResourceAttr(dataSourceName, "region", testRegion),
					resource.TestCheckResourceAttr(dataSourceName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),

					resource.TestCheckResourceAttr(dataSourceName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "allow_list.0", "10.0.0.3/32"),

					resource.TestCheckResourceAttr(dataSourceName, "updates.day_of_week", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.frequency", "weekly"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.hour_of_day", "3"),

					resource.TestCheckResourceAttr(dataSourceName, "pending_updates.#", "0"),

					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_password_encryption",
						"md5",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_stat_monitor_enable",
						"false",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pglookout_max_failover_replication_time_lag",
						"60",
					),
				),
			},
		},
	})
}

func TestAccDataSource_engineConfig(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_database_postgresql_v2.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataEngineConfig(t, tmpl.TemplateDataEngineConfig{
					Label:    label,
					Region:   testRegion,
					EngineID: "postgresql/14",
					Type:     "g6-nanode-1",

					EngineConfigPGAutovacuumAnalyzeScaleFactor:         0.1,
					EngineConfigPGAutovacuumAnalyzeThreshold:           50,
					EngineConfigPGAutovacuumMaxWorkers:                 3,
					EngineConfigPGAutovacuumNaptime:                    100,
					EngineConfigPGAutovacuumVacuumCostDelay:            20,
					EngineConfigPGAutovacuumVacuumCostLimit:            200,
					EngineConfigPGAutovacuumVacuumScaleFactor:          0.2,
					EngineConfigPGAutovacuumVacuumThreshold:            100,
					EngineConfigPGBGWriterDelay:                        1000,
					EngineConfigPGBGWriterFlushAfter:                   512,
					EngineConfigPGBGWriterLRUMaxpages:                  100,
					EngineConfigPGBGWriterLRUMultiplier:                2.5,
					EngineConfigPGDeadlockTimeout:                      1000,
					EngineConfigPGDefaultToastCompression:              "pglz",
					EngineConfigPGIdleInTransactionSessionTimeout:      60000,
					EngineConfigPGJIT:                                  true,
					EngineConfigPGMaxFilesPerProcess:                   1000,
					EngineConfigPGMaxLocksPerTransaction:               64,
					EngineConfigPGMaxLogicalReplicationWorkers:         4,
					EngineConfigPGMaxParallelWorkers:                   8,
					EngineConfigPGMaxParallelWorkersPerGather:          2,
					EngineConfigPGMaxPredLocksPerTransaction:           128,
					EngineConfigPGMaxReplicationSlots:                  8,
					EngineConfigPGMaxSlotWALKeepSize:                   128,
					EngineConfigPGMaxStackDepth:                        2097152,
					EngineConfigPGMaxStandbyArchiveDelay:               60000,
					EngineConfigPGMaxStandbyStreamingDelay:             60000,
					EngineConfigPGMaxWALSenders:                        20,
					EngineConfigPGMaxWorkerProcesses:                   8,
					EngineConfigPGPasswordEncryption:                   "scram-sha-256",
					EngineConfigPGPGPartmanBGWInterval:                 3600,
					EngineConfigPGPGPartmanBGWRole:                     "myrolename",
					EngineConfigPGPGStatMonitorPGSMEnableQueryPlan:     true,
					EngineConfigPGPGStatMonitorPGSMMaxBuckets:          5,
					EngineConfigPGPGStatStatementsTrack:                "all",
					EngineConfigPGTempFileLimit:                        100,
					EngineConfigPGTimezone:                             "Europe/Helsinki",
					EngineConfigPGTrackActivityQuerySize:               2048,
					EngineConfigPGTrackCommitTimestamp:                 "on",
					EngineConfigPGTrackFunctions:                       "all",
					EngineConfigPGTrackIOTiming:                        "on",
					EngineConfigPGWALSenderTimeout:                     60000,
					EngineConfigPGWALWriterDelay:                       200,
					EngineConfigPGStatMonitorEnable:                    true,
					EngineConfigPGLookoutMaxFailoverReplicationTimeLag: 10000,
					EngineConfigSharedBuffersPercentage:                25.5,
					EngineConfigWorkMem:                                400,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),

					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_analyze_scale_factor",
						"0.1",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_analyze_threshold",
						"50",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_max_workers",
						"3",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_naptime",
						"100",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_vacuum_cost_delay",
						"20",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_vacuum_cost_limit",
						"200",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_vacuum_scale_factor",
						"0.2",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_autovacuum_vacuum_threshold",
						"100",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_bgwriter_delay",
						"1000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_bgwriter_flush_after",
						"512",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_bgwriter_lru_maxpages",
						"100",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_bgwriter_lru_multiplier",
						"2.5",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_deadlock_timeout",
						"1000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_default_toast_compression",
						"pglz",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_idle_in_transaction_session_timeout",
						"60000",
					),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_pg_jit", "true"),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_files_per_process",
						"1000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_locks_per_transaction",
						"64",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_logical_replication_workers",
						"4",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_parallel_workers",
						"8",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_parallel_workers_per_gather",
						"2",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_pred_locks_per_transaction",
						"128",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_replication_slots",
						"8",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_slot_wal_keep_size",
						"128",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_stack_depth",
						"2097152",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_standby_archive_delay",
						"60000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_standby_streaming_delay",
						"60000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_wal_senders",
						"20",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_max_worker_processes",
						"8",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_password_encryption",
						"scram-sha-256",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_pg_partman_bgw_interval",
						"3600",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_pg_partman_bgw_role",
						"myrolename",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan",
						"true",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_pg_stat_monitor_pgsm_max_buckets",
						"5",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_pg_stat_statements_track",
						"all",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_temp_file_limit",
						"100",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_timezone",
						"Europe/Helsinki",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_track_activity_query_size",
						"2048",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_track_commit_timestamp",
						"on",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_track_functions",
						"all",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_track_io_timing",
						"on",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_wal_sender_timeout",
						"60000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_wal_writer_delay",
						"200",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pg_stat_monitor_enable",
						"true",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_pglookout_max_failover_replication_time_lag",
						"10000",
					),
					resource.TestCheckResourceAttr(
						dataSourceName,
						"engine_config_shared_buffers_percentage",
						"25.5",
					),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_work_mem", "400"),
				),
			},
		},
	})
}
