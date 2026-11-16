//go:build integration || databasevalkeyv2

package databasevalkeyv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/databasevalkeyv2/tmpl"
)

func TestAccDataSourceDatabaseValkeyV2_basic(t *testing.T) {
	t.Parallel()

	const dataSourceName = "data.linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Data(t, tmpl.TemplateData{
					Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2",
					AllowedIP: "10.0.0.3/32", ClusterSize: 1,
					Updates: tmpl.TemplateDataUpdates{
						HourOfDay: 3, DayOfWeek: 2, Duration: 4, Frequency: "weekly",
					},
				}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_id"), knownvalue.StringExact(testEngine)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("allow_list"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("10.0.0.3/32"),
					})),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("updates").AtMapKey("hour_of_day"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("updates").AtMapKey("day_of_week"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("updates").AtMapKey("duration"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("updates").AtMapKey("frequency"), knownvalue.StringExact("weekly")),
				},
			},
		},
	})
}

func TestAccDataSourceDatabaseValkeyV2_listing(t *testing.T) {
	t.Parallel()
	const databaseName = "data.linode_databases.foobar"
	const enginesName = "data.linode_database_engines.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataListing(t, label, testRegion, testEngine),
				Check:  resource.TestCheckResourceAttrPair(databaseName, "databases.0.id", "linode_database_valkey_v2.foobar", "id"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(databaseName, tfjsonpath.New("databases"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(databaseName, tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("engine"), knownvalue.StringExact("valkey")),
					statecheck.ExpectKnownValue(databaseName, tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(databaseName, tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("oldest_restore_time"), knownvalue.Null()),
					statecheck.ExpectKnownValue(enginesName, tfjsonpath.New("engines"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(enginesName, tfjsonpath.New("engines").AtSliceIndex(0).AtMapKey("engine"), knownvalue.StringExact("valkey")),
					statecheck.ExpectKnownValue(enginesName, tfjsonpath.New("engines").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact(testEngine)),
					statecheck.ExpectKnownValue(enginesName, tfjsonpath.New("engines").AtSliceIndex(0).AtMapKey("version"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDataSourceDatabaseValkeyV2_engineConfig(t *testing.T) {
	t.Parallel()

	const dataSourceName = "data.linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	config := tmpl.TemplateDataEngineConfig{
		Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2",
		EngineConfigBackupHour: 1, EngineConfigBackupMinute: 15, EngineConfigFrequentSnapshots: true,
		EngineConfigValkeyACLChannelsDefault: "allchannels", EngineConfigValkeyActiveExpireEffort: 2,
		EngineConfigValkeyActiveDefrag: true, EngineConfigValkeyLFUDecayTime: 1, EngineConfigValkeyLFULogFactor: 10,
		EngineConfigValkeyMaxmemoryPolicy: "allkeys-lru", EngineConfigValkeyNumberOfDatabases: 16,
		EngineConfigValkeyPersistence: "rdb", EngineConfigValkeyPubsubClientOutputBufferLimit: 32768,
		EngineConfigValkeyTimeout: 300,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataEngineConfig(t, config),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_backup_hour"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_backup_minute"), knownvalue.Int64Exact(15)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_frequent_snapshots"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("engine_config_valkey_acl_channels_default"),
						knownvalue.StringExact("allchannels"),
					),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_active_expire_effort"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_activedefrag"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_lfu_decay_time"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_lfu_log_factor"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_maxmemory_policy"), knownvalue.StringExact("allkeys-lru")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_number_of_databases"), knownvalue.Int64Exact(16)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_persistence"), knownvalue.StringExact("rdb")),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("engine_config_valkey_pubsub_client_output_buffer_limit"),
						knownvalue.Int64Exact(32768),
					),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_timeout"), knownvalue.Int64Exact(300)),
				},
			},
		},
	})
}

func TestAccDataSourceDatabaseValkeyV2_vpc(t *testing.T) {
	t.Parallel()

	const dataSourceName = "data.linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyVPCPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataVPC(t, label, testVPCRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network").AtMapKey("public_access"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network").AtMapKey("vpc_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network").AtMapKey("subnet_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
