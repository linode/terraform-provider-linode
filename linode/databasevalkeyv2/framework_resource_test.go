//go:build integration || databasevalkeyv2

package databasevalkeyv2_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/databasevalkeyv2/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

var testRegion, testVPCRegion, testEngine string
var valkeyAvailabilityError error
var valkeyVPCAvailabilityError error

func init() {
	resource.AddTestSweepers("linode_database_valkey_v2", &resource.Sweeper{
		Name: "linode_database_valkey_v2",
		F:    sweep,
	})

	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(err)
	}

	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityDBAAS}, "core")
	if err != nil {
		valkeyAvailabilityError = err
		return
	}
	testRegion = region

	vpcRegion, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{
		linodego.CapabilityDBAAS,
		linodego.CapabilityVPCs,
	}, "core")
	if err != nil {
		valkeyVPCAvailabilityError = err
	} else {
		testVPCRegion = vpcRegion
	}

	engine, err := databaseshared.ResolveValidDBEngine(context.Background(), *client, string(linodego.DatabaseEngineTypeValkey))
	if err != nil {
		valkeyAvailabilityError = err
		return
	}
	testEngine = engine.ID
}

func valkeyPreCheck(t *testing.T) {
	acceptance.PreCheck(t)
	if valkeyAvailabilityError != nil {
		t.Skipf("Valkey is not available in this test account: %v", valkeyAvailabilityError)
	}
}

func valkeyVPCPreCheck(t *testing.T) {
	valkeyPreCheck(t)
	if valkeyVPCAvailabilityError != nil {
		t.Skipf("Valkey VPC is not available in this test account: %v", valkeyVPCAvailabilityError)
	}
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	databases, err := client.ListValkeyDatabases(context.Background(), acceptance.SweeperListOptions(prefix, "label"))
	if err != nil {
		return fmt.Errorf("failed to list Valkey databases: %w", err)
	}
	for _, database := range databases {
		if !acceptance.ShouldSweep(prefix, database.Label) {
			continue
		}
		if err := client.DeleteValkeyDatabase(context.Background(), database.ID); err != nil {
			return fmt.Errorf("failed to delete Valkey database %s during sweep: %w", database.Label, err)
		}
	}
	return nil
}

func checkValkeyDatabaseDestroy(state *terraform.State) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	for _, instance := range state.RootModule().Resources {
		if instance.Type != "linode_database_valkey_v2" {
			continue
		}
		id, err := strconv.Atoi(instance.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse database ID %q: %w", instance.Primary.ID, err)
		}
		if _, err := client.GetValkeyDatabase(context.Background(), id); err == nil || !linodego.IsNotFound(err) {
			return fmt.Errorf("Valkey database %d still exists or returned an unexpected error: %v", id, err)
		}
	}
	return nil
}

func waitForValkeySnapshot(t *testing.T, label string) {
	t.Helper()

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		databases, err := client.ListValkeyDatabases(ctx, &linodego.ListOptions{})
		if err != nil {
			t.Fatalf("failed to list Valkey databases while waiting for snapshots: %v", err)
		}
		for _, database := range databases {
			if database.Label == label && len(database.AvailableRestoreTimes) > 0 {
				return
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			t.Fatalf("timed out waiting for Valkey database %q to expose a restore snapshot", label)
			return
		}
	}
}

func deleteValkeyDatabaseAndWait(t *testing.T, databaseID int) {
	t.Helper()

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := client.DeleteValkeyDatabase(ctx, databaseID); err != nil {
		t.Fatalf("failed to delete Valkey database %d: %v", databaseID, err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		_, err := client.GetValkeyDatabase(ctx, databaseID)
		if linodego.IsNotFound(err) {
			return
		}
		if err != nil {
			t.Fatalf("failed to verify deletion of Valkey database %d: %v", databaseID, err)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			t.Fatalf("timed out waiting for Valkey database %d to be deleted", databaseID)
			return
		}
	}
}

func TestAccResourceDatabaseValkeyV2_basic(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_id"), knownvalue.StringExact(testEngine)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			{
				Config:       tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				RefreshState: true,
				PlanOnly:     true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_externalDeletion(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	databaseID := 0

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				Check: func(state *terraform.State) error {
					resourceState := state.RootModule().Resources[resourceName]
					id, err := strconv.Atoi(resourceState.Primary.ID)
					if err != nil {
						return fmt.Errorf("failed to parse Valkey database ID %q: %w", resourceState.Primary.ID, err)
					}
					databaseID = id
					return nil
				},
			},
			{
				PreConfig:          func() { deleteValkeyDatabaseAndWait(t, databaseID) },
				Config:             tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_complex(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	updates := tmpl.TemplateData{
		Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2",
		AllowedIP: "10.0.0.3/32", ClusterSize: 1, MaintenanceHour: 3,
	}
	resized := updates
	resized.AllowedIP = "10.0.0.4/32"
	resized.ClusterSize = 3
	cleared := resized
	cleared.AllowedIP = ""
	maintenanceOnly := cleared
	maintenanceOnly.MaintenanceHour = 4

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, updates),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("allow_list"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.3/32")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			{
				Config: tmpl.Complex(t, resized),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("allow_list"), knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.4/32")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			{
				Config: tmpl.Complex(t, cleared),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("allow_list"), knownvalue.SetExact([]knownvalue.Check{})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(3)),
				},
			},
			{
				Config: tmpl.Complex(t, maintenanceOnly),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("allow_list"), knownvalue.SetExact([]knownvalue.Check{})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("updates.hour_of_day"), knownvalue.Int64Exact(4)),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_resize(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(1)),
				},
			},
			{
				Config: tmpl.Resize(t, label, testRegion, testEngine, "g7-dedicated-4-2", 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("cluster_size"), knownvalue.Int64Exact(3)),
				},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_fork(t *testing.T) {
	t.Parallel()

	const sourceName = "linode_database_valkey_v2.foobar"
	const forkName = "linode_database_valkey_v2.fork"
	const latestForkName = "linode_database_valkey_v2.fork_latest"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(sourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			{
				PreConfig: func() { waitForValkeySnapshot(t, label) },
				Config:    tmpl.Fork(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(sourceName, tfjsonpath.New("available_restore_times"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(forkName, tfjsonpath.New("fork_source"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(forkName, tfjsonpath.New("fork_restore_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(forkName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(latestForkName, tfjsonpath.New("fork_source"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(latestForkName, tfjsonpath.New("fork_restore_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(latestForkName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
				},
			},
			{
				ResourceName:            forkName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_suspension(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	data := tmpl.TemplateData{Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2", Suspended: true}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Suspension(t, data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("suspended")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("suspended"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("root_password"), knownvalue.Null()),
				},
			},
			{
				Config: tmpl.Suspension(t, func() tmpl.TemplateData { data.Suspended = false; return data }()),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("suspended"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("root_password"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.Suspension(t, func() tmpl.TemplateData { data.Suspended = true; return data }()),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("suspended")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("suspended"), knownvalue.Bool(true)),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times", "ca_cert", "root_password", "root_username"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_engineConfig(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	data := tmpl.TemplateData{
		Label: label, Region: testRegion, EngineID: testEngine, Type: "g7-dedicated-4-2",
		BackupHour: 1, BackupMinute: 15, FrequentSnapshots: true, ACLChannelsDefault: "allchannels",
		ActiveExpireEffort: 2, ActiveDefrag: true, LFUDecayTime: 1, LFULogFactor: 10,
		MaxmemoryPolicy: "allkeys-lru", NumberOfDatabases: 16, Persistence: "rdb",
		PubsubClientOutputBufferLimit: 32768, Timeout: 300,
	}
	updated := data
	updated.BackupHour = 2
	updated.BackupMinute = 30
	updated.FrequentSnapshots = false
	updated.ACLChannelsDefault = "resetchannels"
	updated.ActiveExpireEffort = 4
	updated.ActiveDefrag = false
	updated.LFUDecayTime = 2
	updated.LFULogFactor = 20
	updated.Timeout = 600
	updated.MaxmemoryPolicy = "volatile-lru"
	updated.NumberOfDatabases = 32
	updated.PubsubClientOutputBufferLimit = 65536

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.EngineConfig(t, data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_backup_hour"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_backup_minute"), knownvalue.Int64Exact(15)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_frequent_snapshots"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_acl_channels_default"), knownvalue.StringExact("allchannels")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_active_expire_effort"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_activedefrag"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_lfu_decay_time"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_lfu_log_factor"), knownvalue.Int64Exact(10)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_timeout"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_maxmemory_policy"), knownvalue.StringExact("allkeys-lru")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_number_of_databases"), knownvalue.Int64Exact(16)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_persistence"), knownvalue.StringExact("rdb")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_pubsub_client_output_buffer_limit"), knownvalue.Int64Exact(32768)),
				},
			},
			{
				Config: tmpl.EngineConfig(t, updated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_backup_hour"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_backup_minute"), knownvalue.Int64Exact(30)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_frequent_snapshots"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_acl_channels_default"), knownvalue.StringExact("resetchannels")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_active_expire_effort"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_activedefrag"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_lfu_decay_time"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_lfu_log_factor"), knownvalue.Int64Exact(20)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_timeout"), knownvalue.Int64Exact(600)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_maxmemory_policy"), knownvalue.StringExact("volatile-lru")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_number_of_databases"), knownvalue.Int64Exact(32)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_persistence"), knownvalue.StringExact("rdb")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("engine_config_valkey_pubsub_client_output_buffer_limit"), knownvalue.Int64Exact(65536)),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_vpc(t *testing.T) {
	t.Parallel()

	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { valkeyVPCPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPC0(t, label, testVPCRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("private_network.public_access"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("private_network.vpc_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("private_network.subnet_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.VPC1(t, label, testVPCRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("private_network.public_access"), knownvalue.Bool(true)),
				},
			},
			{
				Config: tmpl.VPCDetach(t, label, testVPCRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("private_network"), knownvalue.Null()),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}

func TestAccResourceDatabaseValkeyV2_noPendingUpdatesRegression(t *testing.T) {
	t.Parallel()

	overriddenProvider := acceptance.NewFrameworkProviderWithClient(acceptance.NewClientWithDatabasePendingUpdates(t))
	const resourceName = "linode_database_valkey_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { valkeyPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](overriddenProvider, nil)
			},
		},
		CheckDestroy: checkValkeyDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("pending_updates"), acceptance.DatabasePendingUpdatesSetExact),
				},
			},
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g7-dedicated-4-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("pending_updates"), acceptance.DatabasePendingUpdatesSetExact),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "members", "available_restore_times"},
			},
		},
	})
}
