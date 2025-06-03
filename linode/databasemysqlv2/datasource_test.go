//go:build integration || databasemysqlv2

package databasemysqlv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databasemysqlv2/tmpl"
)

func TestAccDataSource_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_database_mysql_v2.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
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
					resource.TestCheckResourceAttr(dataSourceName, "engine", "mysql"),
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

					resource.TestCheckResourceAttr(dataSourceName, "updates.hour_of_day", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.day_of_week", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.frequency", "weekly"),

					resource.TestCheckResourceAttr(dataSourceName, "pending_updates.#", "0"),

					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_sql_mode", "ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_sql_require_primary_key", "true"),
				),
			},
		},
	})
}

func TestAccDataSource_engineConfig(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_database_mysql_v2.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataEngineConfig(t, tmpl.TemplateDataEngineConfig{
					Label:    label,
					Region:   testRegion,
					EngineID: testEngine,
					Type:     "g6-nanode-1",

					EngineConfigBinlogRetentionPeriod:             600,
					EngineConfigMySQLConnectTimeout:               15,
					EngineConfigMySQLDefaultTimeZone:              "+02:00",
					EngineConfigMySQLGroupConcatMaxLen:            2048,
					EngineConfigMySQLInformationSchemaStatsExpiry: 7200,
					EngineConfigMySQLInnoDBChangeBufferMaxSize:    30,
					EngineConfigMySQLInnoDBFlushNeighbors:         0,
					EngineConfigMySQLInnoDBFTMinTokenSize:         4,
					EngineConfigMySQLInnoDBFTServerStopwordTable:  "mysql/innodb_ft_custom_stopword",
					EngineConfigMySQLInnoDBLockWaitTimeout:        600,
					EngineConfigMySQLInnoDBLogBufferSize:          33554432,
					EngineConfigMySQLInnoDBOnlineAlterLogMaxSize:  536870912,
					EngineConfigMySQLInnoDBReadIOThreads:          8,
					EngineConfigMySQLInnoDBRollbackOnTimeout:      false,
					EngineConfigMySQLInnoDBThreadConcurrency:      16,
					EngineConfigMySQLInnoDBWriteIOThreads:         8,
					EngineConfigMySQLInteractiveTimeout:           600,
					EngineConfigMySQLInternalTmpMemStorageEngine:  "TempTable",
					EngineConfigMySQLMaxAllowedPacket:             134217728,
					EngineConfigMySQLMaxHeapTableSize:             33554432,
					EngineConfigMySQLNetBufferLength:              32768,
					EngineConfigMySQLNetReadTimeout:               60,
					EngineConfigMySQLNetWriteTimeout:              60,
					EngineConfigMySQLSortBufferSize:               524288,
					EngineConfigMySQLSQLMode:                      "TRADITIONAL,ANSI",
					EngineConfigMySQLSQLRequirePrimaryKey:         false,
					EngineConfigMySQLTmpTableSize:                 33554432,
					EngineConfigMySQLWaitTimeout:                  36000,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),

					resource.TestCheckResourceAttr(dataSourceName, "engine_config_binlog_retention_period", "600"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_connect_timeout", "15"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_default_time_zone", "+02:00"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_group_concat_max_len", "2048"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_information_schema_stats_expiry", "7200"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_change_buffer_max_size", "30"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_flush_neighbors", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_ft_min_token_size", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_ft_server_stopword_table", "mysql/innodb_ft_custom_stopword"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_lock_wait_timeout", "600"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_log_buffer_size", "33554432"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_online_alter_log_max_size", "536870912"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_read_io_threads", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_rollback_on_timeout", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_thread_concurrency", "16"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_innodb_write_io_threads", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_interactive_timeout", "600"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_internal_tmp_mem_storage_engine", "TempTable"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_max_allowed_packet", "134217728"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_max_heap_table_size", "33554432"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_net_buffer_length", "32768"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_net_read_timeout", "60"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_net_write_timeout", "60"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_sort_buffer_size", "524288"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_sql_mode", "TRADITIONAL,ANSI"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_sql_require_primary_key", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_tmp_table_size", "33554432"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_config_mysql_wait_timeout", "36000"),
				),
			},
		},
	})
}
