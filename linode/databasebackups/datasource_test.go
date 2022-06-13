package databasebackups_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databasebackups"
	"github.com/linode/terraform-provider-linode/linode/databasebackups/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var engineVersion string

func init() {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mongodb")
	if err != nil {
		log.Fatalf("failde to get db engine version: %s", err)
	}

	engineVersion = v.ID
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

func TestFlattenBackup_MongoDB(t *testing.T) {
	currentTime := time.Now()

	backup := linodego.MongoDatabaseBackup{
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

func TestAccDataSourceMongoBackups_basic(t *testing.T) {
	t.Parallel()

	var db linodego.MongoDatabase

	const backupLabel = "coolbackup42"
	dbLabel := acctest.RandomWithPrefix("tf_test")

	resourceName := "linode_database_mongodb.foobar"
	dataSourceName := "data.linode_database_backups.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, tmpl.TemplateData{
					Engine:      engineVersion,
					Label:       dbLabel,
					BackupLabel: backupLabel,
				}),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMongoDatabaseExists(resourceName, &db),
					resource.TestCheckResourceAttr(dataSourceName, "backups.#", "0"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

					if err := client.CreateMongoDatabaseBackup(context.Background(), db.ID, linodego.MongoBackupCreateOptions{
						Label:  backupLabel,
						Target: "primary",
					}); err != nil {
						t.Errorf("failed to create db backup: %v", err)
					}

					err := client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypeMongo, "backing_up", 120)
					if err != nil {
						t.Fatalf("failed to wait for database backing_up: %s", err)
					}

					err = client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypeMongo, linodego.DatabaseStatusActive, 1200)
					if err != nil {
						t.Fatalf("failed to wait for database active: %s", err)
					}
				},
				Config: tmpl.DataBasic(t, tmpl.TemplateData{
					Engine:      engineVersion,
					Label:       dbLabel,
					BackupLabel: backupLabel,
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
