package databasemysqlbackups_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databasemysqlbackups/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var engineVersion string

func init() {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failde to get db engine version: %s", err)
	}

	engineVersion = v.ID
}

func TestAccDataSourceMySQLBackups_basic(t *testing.T) {
	t.Parallel()

	var db linodego.MySQLDatabase

	const backupLabel = "coolbackup42"
	dbLabel := acctest.RandomWithPrefix("tf_test")

	resourceName := "linode_database_mysql.foobar"
	dataSourceName := "data.linode_database_mysql_backups.foobar"

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
					acceptance.CheckMySQLDatabaseExists(resourceName, &db),
					resource.TestCheckResourceAttr(dataSourceName, "backups.#", "0"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

					if err := client.CreateMySQLDatabaseBackup(context.Background(), db.ID, linodego.MySQLBackupCreateOptions{
						Label:  backupLabel,
						Target: "primary",
					}); err != nil {
						t.Errorf("failed to create db backup: %v", err)
					}

					err := client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypeMySQL, "backing_up", 120)
					if err != nil {
						t.Fatalf("failed to wait for database backing_up: %s", err)
					}

					err = client.WaitForDatabaseStatus(context.Background(), db.ID,
						linodego.DatabaseEngineTypeMySQL, linodego.DatabaseStatusActive, 1200)
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
