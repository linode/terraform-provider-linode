//go:build integration || databasepostgresqlv2

package databasepostgresqlv2_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databasepostgresqlv2/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region

	engine, err := helper.ResolveValidDBEngine(
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
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

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),

					resource.TestCheckResourceAttrSet(resName, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resName, "updates.duration"),
					resource.TestCheckResourceAttrSet(resName, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resName, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resName, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
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
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
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
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
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
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
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

					resource.TestCheckResourceAttr(resNameSource, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resNameSource, "allow_list.0", "0.0.0.0/0"),

					resource.TestCheckResourceAttrSet(resNameSource, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.duration"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resNameSource, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resNameSource, "pending_updates.#", "0"),
				),
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

					resource.TestCheckResourceAttr(resNameFork, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resNameFork, "allow_list.0", "0.0.0.0/0"),

					resource.TestCheckResourceAttrSet(resNameFork, "updates.day_of_week"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.duration"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.frequency"),
					resource.TestCheckResourceAttrSet(resNameFork, "updates.hour_of_day"),

					resource.TestCheckResourceAttr(resNameFork, "pending_updates.#", "0"),
				),
			},
			{
				ResourceName:      resNameSource,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resNameFork,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
