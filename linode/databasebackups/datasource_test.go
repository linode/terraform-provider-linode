//go:build integration || databasebackups

package databasebackups_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/databasebackups/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

var (
	engineVersion string
	testRegion    string
)

func init() {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "postgresql")
	if err != nil {
		log.Fatalf("failde to get db engine version: %s", err)
	}

	engineVersion = v.ID

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourcePostgresBackups_basic(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	var db linodego.PostgresDatabase

	const backupLabel = "coolbackup42"
	dbLabel := acctest.RandomWithPrefix("tf_test")

	resourceName := "linode_database_postgresql.foobar"
	dataSourceName := "data.linode_database_backups.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, tmpl.TemplateData{
					Engine:      engineVersion,
					Region:      testRegion,
					Label:       dbLabel,
					BackupLabel: backupLabel,
				}),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckPostgresDatabaseExists(resourceName, &db),
					resource.TestCheckResourceAttr(dataSourceName, "backups.#", "0"),
				),
			},
			{
				PreConfig: func() {
					client, err := acceptance.GetTestClient()
					if err != nil {
						log.Fatalf("failed to get client: %s", err)
					}

					if err := client.CreatePostgresDatabaseBackup(context.Background(), db.ID, linodego.PostgresBackupCreateOptions{
						Label:  backupLabel,
						Target: "primary",
					}); err != nil {
						t.Errorf("failed to create db backup: %v", err)
					}

					err = client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypePostgres, "backing_up", 120)
					if err != nil {
						t.Fatalf("failed to wait for database backing_up: %s", err)
					}

					err = client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypePostgres, linodego.DatabaseStatusActive, 1200)
					if err != nil {
						t.Fatalf("failed to wait for database active: %s", err)
					}

					for {
						list, err := client.ListPostgresDatabaseBackups(context.Background(), db.ID, nil)
						if err != nil {
							t.Fatalf("failed to list database backups: %s", err)
						}

						if len(list) > 0 {
							break
						}
					}
				},
				Config: tmpl.DataBasic(t, tmpl.TemplateData{
					Engine:      engineVersion,
					Label:       dbLabel,
					BackupLabel: backupLabel,
					Region:      testRegion,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "backups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "backups.0.label", backupLabel),
					resource.TestCheckResourceAttr(dataSourceName, "backups.0.type", "snapshot"),
				),
			},
		},
	})
}
