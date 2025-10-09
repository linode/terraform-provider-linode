//go:build integration || databasemysqlv2

package databasemysqlv2_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/databasemysqlv2/tmpl"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

var testRegion, testEngine string

func init() {
	resource.AddTestSweepers("linode_database_mysql_v2", &resource.Sweeper{
		Name: "linode_database_mysql_v2",
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

	engine, err := helper.ResolveValidDBEngine(
		context.Background(),
		*client,
		string(linodego.DatabaseEngineTypeMySQL),
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

	dbs, err := client.ListMySQLDatabases(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("error getting mysql databases: %w", err)
	}
	for _, db := range dbs {
		if !acceptance.ShouldSweep(prefix, db.Label) {
			continue
		}
		err := client.DeleteMySQLDatabase(context.Background(), db.ID)
		if err != nil {
			return fmt.Errorf("error destroying %s during sweep: %w", db.Label, err)
		}
	}

	return nil
}

func TestAccResource_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "mysql"),
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

					resource.TestCheckResourceAttr(
						resName,
						"engine_config_mysql_sql_mode",
						"ANSI,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,STRICT_ALL_TABLES",
					),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_require_primary_key", "true"),
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

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "mysql"),
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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "mysql"),
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

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "mysql"),
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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "3"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "engine", "mysql"),
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

	resNameSource := "linode_database_mysql_v2.foobar"
	resNameFork := "linode_database_mysql_v2.fork"

	var dbSource linodego.MySQLDatabase

	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resNameSource, &dbSource),

					resource.TestCheckResourceAttrSet(resNameSource, "id"),

					resource.TestCheckResourceAttrSet(resNameSource, "ca_cert"),
					resource.TestCheckResourceAttr(resNameSource, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resNameSource, "created"),
					resource.TestCheckResourceAttr(resNameSource, "encrypted", "true"),
					resource.TestCheckResourceAttr(resNameSource, "engine", "mysql"),
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
							db, err := client.GetMySQLDatabase(ctx, dbSource.ID)
							if err != nil {
								t.Fatalf("failed to get mysql database: %s", err)
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
					acceptance.CheckMySQLDatabaseExists(resNameFork, nil),

					resource.TestCheckResourceAttrSet(resNameFork, "id"),

					resource.TestCheckResourceAttrSet(resNameFork, "ca_cert"),
					resource.TestCheckResourceAttr(resNameFork, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(resNameFork, "created"),
					resource.TestCheckResourceAttr(resNameFork, "encrypted", "true"),
					resource.TestCheckResourceAttr(resNameFork, "engine", "mysql"),
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

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

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

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.EngineConfig(
					t,
					tmpl.TemplateDataEngineConfig{
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
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "engine_config_binlog_retention_period", "600"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_connect_timeout", "15"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_default_time_zone", "+02:00"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_group_concat_max_len", "2048"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_information_schema_stats_expiry", "7200"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_change_buffer_max_size", "30"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_flush_neighbors", "0"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_ft_min_token_size", "4"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_ft_server_stopword_table", "mysql/innodb_ft_custom_stopword"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_lock_wait_timeout", "600"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_log_buffer_size", "33554432"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_online_alter_log_max_size", "536870912"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_read_io_threads", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_rollback_on_timeout", "false"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_thread_concurrency", "16"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_write_io_threads", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_interactive_timeout", "600"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_internal_tmp_mem_storage_engine", "TempTable"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_allowed_packet", "134217728"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_heap_table_size", "33554432"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_buffer_length", "32768"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_read_timeout", "60"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_write_timeout", "60"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sort_buffer_size", "524288"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_mode", "TRADITIONAL,ANSI"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_require_primary_key", "false"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_tmp_table_size", "33554432"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_wait_timeout", "36000"),
				),
			},
			{
				Config: tmpl.EngineConfig(
					t,
					tmpl.TemplateDataEngineConfig{
						Label:    label,
						Region:   testRegion,
						EngineID: testEngine,
						Type:     "g6-nanode-1",

						EngineConfigBinlogRetentionPeriod:             1200,
						EngineConfigMySQLConnectTimeout:               30,
						EngineConfigMySQLDefaultTimeZone:              "+05:00",
						EngineConfigMySQLGroupConcatMaxLen:            4096,
						EngineConfigMySQLInformationSchemaStatsExpiry: 18000,
						EngineConfigMySQLInnoDBChangeBufferMaxSize:    15,
						EngineConfigMySQLInnoDBFlushNeighbors:         1,
						EngineConfigMySQLInnoDBFTMinTokenSize:         8,
						EngineConfigMySQLInnoDBFTServerStopwordTable:  "db_name/innodb_ft_stopword_list",
						EngineConfigMySQLInnoDBLockWaitTimeout:        300,
						EngineConfigMySQLInnoDBLogBufferSize:          67108864,
						EngineConfigMySQLInnoDBOnlineAlterLogMaxSize:  1342177280,
						EngineConfigMySQLInnoDBReadIOThreads:          6,
						EngineConfigMySQLInnoDBRollbackOnTimeout:      true,
						EngineConfigMySQLInnoDBThreadConcurrency:      20,
						EngineConfigMySQLInnoDBWriteIOThreads:         10,
						EngineConfigMySQLInteractiveTimeout:           900,
						EngineConfigMySQLInternalTmpMemStorageEngine:  "MEMORY",
						EngineConfigMySQLMaxAllowedPacket:             134217728,
						EngineConfigMySQLMaxHeapTableSize:             67108864,
						EngineConfigMySQLNetBufferLength:              32768,
						EngineConfigMySQLNetReadTimeout:               90,
						EngineConfigMySQLNetWriteTimeout:              90,
						EngineConfigMySQLSortBufferSize:               1048576,
						EngineConfigMySQLSQLMode:                      "STRICT_TRANS_TABLES,ANSI",
						EngineConfigMySQLSQLRequirePrimaryKey:         true,
						EngineConfigMySQLTmpTableSize:                 67108864,
						EngineConfigMySQLWaitTimeout:                  43200,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),

					resource.TestCheckResourceAttr(resName, "engine_config_binlog_retention_period", "1200"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_connect_timeout", "30"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_default_time_zone", "+05:00"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_group_concat_max_len", "4096"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_information_schema_stats_expiry", "18000"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_change_buffer_max_size", "15"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_flush_neighbors", "1"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_ft_min_token_size", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_ft_server_stopword_table", "db_name/innodb_ft_stopword_list"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_lock_wait_timeout", "300"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_log_buffer_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_online_alter_log_max_size", "1342177280"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_read_io_threads", "6"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_rollback_on_timeout", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_thread_concurrency", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_write_io_threads", "10"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_interactive_timeout", "900"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_internal_tmp_mem_storage_engine", "MEMORY"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_allowed_packet", "134217728"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_heap_table_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_buffer_length", "32768"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_read_timeout", "90"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_write_timeout", "90"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sort_buffer_size", "1048576"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_tmp_table_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_wait_timeout", "43200"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_mode", "STRICT_TRANS_TABLES,ANSI"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_require_primary_key", "true"),
				),
			},
			// Verify that omitting the nullable field EngineConfigMySQLInnoDBFTServerStopwordTable does not affect other engine config fields in the Terraform output
			{
				Config: tmpl.EngineConfigNullableField(
					t,
					tmpl.TemplateDataEngineConfig{
						Label:                             label,
						Region:                            testRegion,
						EngineID:                          testEngine,
						Type:                              "g6-nanode-1",
						EngineConfigBinlogRetentionPeriod: 1800,
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),

					resource.TestCheckResourceAttrSet(resName, "id"),
					// Updated EngineConfigBinlogRetentionPeriod field assertion
					resource.TestCheckResourceAttr(resName, "engine_config_binlog_retention_period", "1800"),
					// Retained fields
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_connect_timeout", "30"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_default_time_zone", "+05:00"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_group_concat_max_len", "4096"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_information_schema_stats_expiry", "18000"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_change_buffer_max_size", "15"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_flush_neighbors", "1"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_ft_min_token_size", "8"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_lock_wait_timeout", "300"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_log_buffer_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_online_alter_log_max_size", "1342177280"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_read_io_threads", "6"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_rollback_on_timeout", "true"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_thread_concurrency", "20"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_innodb_write_io_threads", "10"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_interactive_timeout", "900"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_internal_tmp_mem_storage_engine", "MEMORY"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_allowed_packet", "134217728"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_max_heap_table_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_buffer_length", "32768"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_read_timeout", "90"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_net_write_timeout", "90"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sort_buffer_size", "1048576"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_tmp_table_size", "67108864"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_wait_timeout", "43200"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_mode", "STRICT_TRANS_TABLES,ANSI"),
					resource.TestCheckResourceAttr(resName, "engine_config_mysql_sql_require_primary_key", "true"),
					// Nullable field EngineConfigMySQLInnoDBFTServerStopwordTable assertion
					resource.TestCheckNoResourceAttr(resName, "engine_config_mysql_innodb_ft_server_stopword_table"),
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

	resName := "linode_database_mysql_v2.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPC0(t, label, testRegion, testEngine, "g6-nanode-1"),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),

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
					acceptance.CheckMySQLDatabaseExists(resName, nil),

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
