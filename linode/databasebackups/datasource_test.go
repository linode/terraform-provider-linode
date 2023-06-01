package databasebackups_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/databasebackups/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databasebackups"
)

var (
	engineVersion string
	testRegion    string
)

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestFlattenBackup_MySQL(t *testing.T) {
	currentTime := time.Now()

	backup := linodego.MySQLDatabaseBackup{
		ID:      123,
		Label:   "cool",
		Type:    "auto",
		Created: &currentTime,
	}

	result := databasebackups.FlattenBackup(backup)
	if result["id"] != backup.ID {
		t.Fatal(cmp.Diff(result["id"], backup.ID))
	}

	if result["label"] != backup.Label {
		t.Fatal(cmp.Diff(result["label"], backup.Label))
	}

	if result["type"] != backup.Type {
		t.Fatal(cmp.Diff(result["type"], backup.Type))
	}
}

func TestFlattenBackup_PostgreSQL(t *testing.T) {
	currentTime := time.Now()

	backup := linodego.PostgresDatabaseBackup{
		ID:      123,
		Label:   "cool",
		Type:    "auto",
		Created: &currentTime,
	}

	result := databasebackups.FlattenBackup(backup)
	if result["id"] != backup.ID {
		t.Fatal(cmp.Diff(result["id"], backup.ID))
	}

	if result["label"] != backup.Label {
		t.Fatal(cmp.Diff(result["label"], backup.Label))
	}

	if result["type"] != backup.Type {
		t.Fatal(cmp.Diff(result["type"], backup.Type))
	}
}

func TestAccDataSourcePostgresBackups_basic(t *testing.T) {
	t.Parallel()
	t.Skip()

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
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

					if err := client.CreatePostgresDatabaseBackup(context.Background(), db.ID, linodego.PostgresBackupCreateOptions{
						Label:  backupLabel,
						Target: "primary",
					}); err != nil {
						t.Errorf("failed to create db backup: %v", err)
					}

					err := client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypePostgres, "backing_up", 120)
					if err != nil {
						t.Fatalf("failed to wait for database backing_up: %s", err)
					}

					err = client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypePostgres, linodego.DatabaseStatusActive, 1200)
					if err != nil {
						t.Fatalf("failed to wait for database active: %s", err)
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
