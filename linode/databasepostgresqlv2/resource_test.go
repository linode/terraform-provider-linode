//go:build integration || databasepostgresqlv2

package databasepostgresqlv2_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/databasepostgresqlv2/tmpl"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/databaseshared"
)

var testRegion, testEngine string

func init() {
	resource.AddTestSweepers("linode_database_postgresql_v2", &resource.Sweeper{
		Name: "linode_database_postgresql_v2",
		F:    sweep,
	})

	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatal(err)
	}

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases", "VPCs"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region

	engine, err := databaseshared.ResolveValidDBEngine(
		context.Background(),
		*client,
		string(linodego.DatabaseEngineTypePostgres),
	)
	if err != nil {
		log.Fatal(err)
	}

	testEngine = engine.ID
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")

	dbs, err := client.ListPostgresDatabases(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("error getting postgres databases: %w", err)
	}
	for _, db := range dbs {
		if !acceptance.ShouldSweep(prefix, db.Label) {
			continue
		}
		err := client.DeletePostgresDatabase(context.Background(), db.ID)
		if err != nil {
			return fmt.Errorf("error destroying %s during sweep: %w", db.Label, err)
		}
	}

	return nil
}

func TestAccResource_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resName, "fork_source"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttrSet(resName, "members.%"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resName, "port"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "version"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "0"),

					resource.TestCheckResourceAttrSet(resName, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resName, "updates.duration"),
					resource.TestCheckResourceAttrSet(resName, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resName, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),

					resource.TestCheckResourceAttr(resName, "engine_config_pg_password_encryption", "md5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_stat_monitor_enable", "false"),
					resource.TestCheckResourceAttr(resName, "engine_config_pglookout_max_failover_replication_time_lag", "60"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_resize(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(
					t,
					tmpl.TemplateData{
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
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resName, "fork_source"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttrSet(resName, "members.%"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resName, "port"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "version"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "10.0.0.3/32"),

					resource.TestCheckResourceAttr(resName, "updates.day_of_week", "2"),
					resource.TestCheckResourceAttr(resName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(resName, "updates.frequency", "weekly"),
					resource.TestCheckResourceAttr(resName, "updates.hour_of_day", "3"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),
				),
			},
			{
				Config: tmpl.Complex(
					t,
					tmpl.TemplateData{
						Label:       label,
						Region:      testRegion,
						EngineID:    testEngine,
						Type:        "g6-standard-1",
						AllowedIP:   "10.0.0.3/32",
						ClusterSize: 1,
						Updates: tmpl.TemplateDataUpdates{
							HourOfDay: 3,
							DayOfWeek: 2,
							Duration:  4,
							Frequency: "weekly",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resName, "fork_source"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttrSet(resName, "members.%"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resName, "port"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "version"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "10.0.0.3/32"),

					resource.TestCheckResourceAttr(resName, "updates.day_of_week", "2"),
					resource.TestCheckResourceAttr(resName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(resName, "updates.frequency", "weekly"),
					resource.TestCheckResourceAttr(resName, "updates.hour_of_day", "3"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(
					t,
					tmpl.TemplateData{
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
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resName, "fork_source"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttrSet(resName, "members.%"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resName, "port"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "version"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "10.0.0.3/32"),

					resource.TestCheckResourceAttr(resName, "updates.day_of_week", "2"),
					resource.TestCheckResourceAttr(resName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(resName, "updates.frequency", "weekly"),
					resource.TestCheckResourceAttr(resName, "updates.hour_of_day", "3"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),
				),
			},
			{
				Config: tmpl.Complex(
					t,
					tmpl.TemplateData{
						Label:       label,
						Region:      testRegion,
						EngineID:    testEngine,
						Type:        "g6-nanode-1",
						AllowedIP:   "10.0.0.4/32",
						ClusterSize: 3,
						Updates: tmpl.TemplateDataUpdates{
							HourOfDay: 2,
							DayOfWeek: 3,
							Duration:  4,
							Frequency: "weekly",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "3"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resName, "fork_source"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttrSet(resName, "members.%"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resName, "port"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "version"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "10.0.0.4/32"),

					resource.TestCheckResourceAttr(resName, "updates.hour_of_day", "2"),
					resource.TestCheckResourceAttr(resName, "updates.day_of_week", "3"),
					resource.TestCheckResourceAttr(resName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(resName, "updates.frequency", "weekly"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_fork(t *testing.T) {
	t.Parallel()

	resNameSource := "linode_database_postgresql_v2.foobar"
	resNameFork := "linode_database_postgresql_v2.fork"

	var dbSource linodego.PostgresDatabase

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resNameSource, &dbSource),

					resource.TestCheckResourceAttrSet(resNameSource, "id"),

					resource.TestCheckResourceAttrSet(resNameSource, "ca_cert"),
					resource.TestCheckResourceAttr(resNameSource, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resNameSource, "created"),
					resource.TestCheckResourceAttr(resNameSource, "encrypted", "true"),
					resource.TestCheckResourceAttr(resNameSource, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resNameSource, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(resNameSource, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(resNameSource, "fork_source"),
					resource.TestCheckResourceAttrSet(resNameSource, "host_primary"),
					resource.TestCheckResourceAttr(resNameSource, "label", label),
					resource.TestCheckResourceAttrSet(resNameSource, "members.%"),
					resource.TestCheckResourceAttrSet(resNameSource, "root_password"),
					resource.TestCheckResourceAttrSet(resNameSource, "root_username"),
					resource.TestCheckResourceAttr(resNameSource, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resNameSource, "port"),
					resource.TestCheckResourceAttr(resNameSource, "region", testRegion),
					resource.TestCheckResourceAttr(resNameSource, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resNameSource, "status", "active"),
					resource.TestCheckResourceAttr(resNameSource, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resNameSource, "updated"),
					resource.TestCheckResourceAttrSet(resNameSource, "version"),

					resource.TestCheckResourceAttr(resNameSource, "allow_list.#", "0"),

					resource.TestCheckResourceAttrSet(resNameSource, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.duration"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resNameSource, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:            resNameSource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
			{
				PreConfig: func() {
					// Poll for the source database to be restorable
					ctx := context.Background()

					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatal(err)
					}

					ctx, cancel := context.WithTimeout(ctx, time.Minute*30)
					defer cancel()

					ticker := time.NewTicker(5 * time.Second)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							db, err := client.GetPostgresDatabase(ctx, dbSource.ID)
							if err != nil {
								t.Fatalf("failed to get postgres database: %s", err)
							}

							if db.OldestRestoreTime != nil {
								return
							}
						case <-ctx.Done():
							return
						}
					}
				},
				Config: tmpl.Fork(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resNameFork, nil),

					resource.TestCheckResourceAttrSet(resNameFork, "id"),

					resource.TestCheckResourceAttrSet(resNameFork, "ca_cert"),
					resource.TestCheckResourceAttr(resNameFork, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resNameFork, "created"),
					resource.TestCheckResourceAttr(resNameFork, "encrypted", "true"),
					resource.TestCheckResourceAttr(resNameFork, "engine", "postgresql"),
					resource.TestCheckResourceAttr(resNameFork, "engine_id", testEngine),
					resource.TestCheckResourceAttrSet(resNameFork, "fork_restore_time"),
					resource.TestCheckResourceAttrSet(resNameFork, "fork_source"),
					resource.TestCheckResourceAttrSet(resNameFork, "host_primary"),
					resource.TestCheckResourceAttr(resNameFork, "label", label+"-fork"),
					resource.TestCheckResourceAttrSet(resNameFork, "members.%"),
					resource.TestCheckResourceAttrSet(resNameSource, "oldest_restore_time"),
					resource.TestCheckResourceAttr(resNameFork, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(resNameFork, "port"),
					resource.TestCheckResourceAttr(resNameFork, "region", testRegion),
					resource.TestCheckResourceAttrSet(resNameFork, "root_password"),
					resource.TestCheckResourceAttrSet(resNameFork, "root_username"),
					resource.TestCheckResourceAttr(resNameFork, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(resNameFork, "status", "active"),
					resource.TestCheckResourceAttr(resNameFork, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(resNameFork, "updated"),
					resource.TestCheckResourceAttrSet(resNameFork, "version"),

					resource.TestCheckResourceAttr(resNameFork, "allow_list.#", "0"),

					resource.TestCheckResourceAttrSet(resNameFork, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.duration"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resNameFork, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:            resNameFork,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_suspension(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Suspension(
					t,
					tmpl.TemplateData{
						Label:     label,
						Region:    testRegion,
						EngineID:  testEngine,
						Type:      "g6-nanode-1",
						Suspended: true,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckNoResourceAttr(resName, "ca_cert"),
					resource.TestCheckNoResourceAttr(resName, "root_password"),
					resource.TestCheckNoResourceAttr(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "status", "suspended"),
					resource.TestCheckResourceAttr(resName, "suspended", "true"),
				),
			},
			{
				Config: tmpl.Suspension(
					t,
					tmpl.TemplateData{
						Label:     label,
						Region:    testRegion,
						EngineID:  testEngine,
						Type:      "g6-nanode-1",
						Suspended: false,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttr(resName, "suspended", "false"),
				),
			},
			{
				Config: tmpl.Suspension(
					t,
					tmpl.TemplateData{
						Label:     label,
						Region:    testRegion,
						EngineID:  testEngine,
						Type:      "g6-nanode-1",
						Suspended: true,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttrSet(resName, "root_username"),
					resource.TestCheckResourceAttr(resName, "status", "suspended"),
					resource.TestCheckResourceAttr(resName, "suspended", "true"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"updated",
					"oldest_restore_time",
					"members",

					// These fields will be populated with null when importing a suspended database
					"ca_cert", "root_password", "root_username",
				},
			},
		},
	})
}

func TestAccResource_engineConfig(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.EngineConfig(
					t,
					tmpl.TemplateDataEngineConfig{
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
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_scale_factor", "0.1"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_threshold", "50"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_max_workers", "3"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_naptime", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_delay", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_limit", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_scale_factor", "0.2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_threshold", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_delay", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_flush_after", "512"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_maxpages", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_multiplier", "2.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_deadlock_timeout", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_default_toast_compression", "pglz"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_idle_in_transaction_session_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_jit", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_files_per_process", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_locks_per_transaction", "64"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_logical_replication_workers", "4"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers_per_gather", "2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_pred_locks_per_transaction", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_replication_slots", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_slot_wal_keep_size", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_stack_depth", "2097152"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_archive_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_streaming_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_wal_senders", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_worker_processes", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_password_encryption", "scram-sha-256"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_interval", "3600"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_role", "myrolename"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_max_buckets", "5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_statements_track", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_temp_file_limit", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_timezone", "Europe/Helsinki"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_activity_query_size", "2048"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_commit_timestamp", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_functions", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_io_timing", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_sender_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_writer_delay", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_stat_monitor_enable", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pglookout_max_failover_replication_time_lag", "10000"),
					resource.TestCheckResourceAttr(resName, "engine_config_shared_buffers_percentage", "25.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_work_mem", "400"),
				),
			},
			{
				Config: tmpl.EngineConfig(
					t,
					tmpl.TemplateDataEngineConfig{
						Label:    label,
						Region:   testRegion,
						EngineID: "postgresql/14",
						Type:     "g6-nanode-1",

						EngineConfigPGAutovacuumAnalyzeScaleFactor:         0.5,
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
						EngineConfigPGPasswordEncryption:                   "md5",
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
						EngineConfigPGLookoutMaxFailoverReplicationTimeLag: 100000,
						EngineConfigSharedBuffersPercentage:                25.5,
						EngineConfigWorkMem:                                400,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_scale_factor", "0.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_threshold", "50"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_max_workers", "3"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_naptime", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_delay", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_limit", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_scale_factor", "0.2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_threshold", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_delay", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_flush_after", "512"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_maxpages", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_multiplier", "2.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_deadlock_timeout", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_default_toast_compression", "pglz"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_idle_in_transaction_session_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_jit", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_files_per_process", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_locks_per_transaction", "64"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_logical_replication_workers", "4"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers_per_gather", "2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_pred_locks_per_transaction", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_replication_slots", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_slot_wal_keep_size", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_stack_depth", "2097152"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_archive_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_streaming_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_wal_senders", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_worker_processes", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_password_encryption", "md5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_interval", "3600"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_role", "myrolename"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_max_buckets", "5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_statements_track", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_temp_file_limit", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_timezone", "Europe/Helsinki"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_activity_query_size", "2048"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_commit_timestamp", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_functions", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_io_timing", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_sender_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_writer_delay", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_stat_monitor_enable", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pglookout_max_failover_replication_time_lag", "100000"),
					resource.TestCheckResourceAttr(resName, "engine_config_shared_buffers_percentage", "25.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_work_mem", "400"),
				),
			},
			// Verify that ommitting or skipping these fields leaves previously set values unchanged
			{
				Config: tmpl.EngineConfigUpdate(
					t,
					tmpl.TemplateDataEngineConfig{
						Label:    label,
						Region:   testRegion,
						EngineID: "postgresql/14",
						Type:     "g6-nanode-1",

						EngineConfigPGAutovacuumAnalyzeScaleFactor: 0.7,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),
					// Updated EngineConfigPGAutovacuumAnalyzeScaleFactor Field
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_scale_factor", "0.7"),
					// Retained Fields
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_analyze_threshold", "50"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_max_workers", "3"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_naptime", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_delay", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_cost_limit", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_scale_factor", "0.2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_autovacuum_vacuum_threshold", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_delay", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_flush_after", "512"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_maxpages", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_bgwriter_lru_multiplier", "2.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_deadlock_timeout", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_default_toast_compression", "pglz"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_idle_in_transaction_session_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_jit", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_files_per_process", "1000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_locks_per_transaction", "64"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_logical_replication_workers", "4"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_parallel_workers_per_gather", "2"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_pred_locks_per_transaction", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_replication_slots", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_slot_wal_keep_size", "128"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_stack_depth", "2097152"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_archive_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_standby_streaming_delay", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_wal_senders", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_max_worker_processes", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_password_encryption", "md5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_interval", "3600"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_partman_bgw_role", "myrolename"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_enable_query_plan", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_monitor_pgsm_max_buckets", "5"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_pg_stat_statements_track", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_temp_file_limit", "100"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_timezone", "Europe/Helsinki"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_activity_query_size", "2048"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_commit_timestamp", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_functions", "all"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_track_io_timing", "on"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_sender_timeout", "60000"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_wal_writer_delay", "200"),
					resource.TestCheckResourceAttr(resName, "engine_config_pg_stat_monitor_enable", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_pglookout_max_failover_replication_time_lag", "100000"),
					resource.TestCheckResourceAttr(resName, "engine_config_shared_buffers_percentage", "25.5"),
					resource.TestCheckResourceAttr(resName, "engine_config_work_mem", "400"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_vpc(t *testing.T) {
	t.Parallel()

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPC0(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "private_network.public_access", "false"),
					resource.TestCheckResourceAttrPair(
						resName, "private_network.vpc_id",
						"linode_vpc.foobar", "id",
					),
					resource.TestCheckResourceAttrPair(
						resName, "private_network.subnet_id",
						"linode_vpc_subnet.foobar", "id",
					),
				),
			},
			{
				Config: tmpl.VPC1(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "private_network.public_access", "true"),
					resource.TestCheckResourceAttrPair(
						resName, "private_network.vpc_id",
						"linode_vpc.foobar2", "id",
					),
					resource.TestCheckResourceAttrPair(
						resName, "private_network.subnet_id",
						"linode_vpc_subnet.foobar2", "id",
					),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}

func TestAccResource_noPendingUpdatesRegression(t *testing.T) {
	t.Parallel()

	overriddenProvider := acceptance.NewFrameworkProviderWithClient(
		acceptance.NewClientWithDatabasePendingUpdates(t),
	)

	resName := "linode_database_postgresql_v2.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](overriddenProvider, nil)
			},
		},
		CheckDestroy: acceptance.CheckPostgreSQLDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("pending_updates"),
						acceptance.DatabasePendingUpdatesSetExact,
					),
				},
			},
			{
				// Ensure refreshes work as expected
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName,
						tfjsonpath.New("pending_updates"),
						acceptance.DatabasePendingUpdatesSetExact,
					),
				},
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated", "oldest_restore_time", "members"},
			},
		},
	})
}
