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
				Config: tmpl.DataBasic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_id"), knownvalue.StringExact(testEngine)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
		},
	})
}

func TestAccDataSourceDatabaseValkeyV2_engineConfig(t *testing.T) {
	t.Parallel()

	const dataSourceName = "data.linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	config := tmpl.TemplateData{
		Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2",
		BackupHour: 1, BackupMinute: 15, FrequentSnapshots: true, ACLChannelsDefault: "allchannels",
		ActiveExpireEffort: 2, ActiveDefrag: true, LFUDecayTime: 1, LFULogFactor: 10,
		MaxmemoryPolicy: "allkeys-lru", NumberOfDatabases: 16, Persistence: "rdb",
		PubsubClientOutputBufferLimit: 32768, Timeout: 300,
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
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_acl_channels_default"), knownvalue.StringExact("allchannels")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_active_expire_effort"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_activedefrag"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_lfu_decay_time"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_lfu_log_factor"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_maxmemory_policy"), knownvalue.StringExact("allkeys-lru")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_number_of_databases"), knownvalue.Int64Exact(16)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_persistence"), knownvalue.StringExact("rdb")),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("engine_config_valkey_pubsub_client_output_buffer_limit"), knownvalue.Int64Exact(32768)),
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
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network.public_access"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network.vpc_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataSourceName, tfjsonpath.New("private_network.subnet_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
