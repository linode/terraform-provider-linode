package databasemongo_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databasemysql/tmpl"
)

func TestAccDataSourceDatabaseMySQL_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_database_mysql.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, tmpl.TemplateData{
					Engine:          engineVersion,
					Label:           dbName,
					AllowedIP:       "0.0.0.0/0",
					ClusterSize:     1,
					Encrypted:       true,
					ReplicationType: "none",
					SSLConnection:   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckMySQLDatabaseExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "engine_id", engineVersion),
					resource.TestCheckResourceAttr(resName, "label", dbName),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),

					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "replication_type", "none"),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),

					resource.TestCheckResourceAttr(resName, "updates.#", "1"),
					resource.TestCheckResourceAttr(resName, "updates.0.day_of_week", "saturday"),
					resource.TestCheckResourceAttr(resName, "updates.0.duration", "1"),
					resource.TestCheckResourceAttr(resName, "updates.0.frequency", "monthly"),
					resource.TestCheckResourceAttr(resName, "updates.0.hour_of_day", "22"),
					resource.TestCheckResourceAttr(resName, "updates.0.week_of_month", "2"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttrSet(resName, "host_secondary"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),

					resource.TestCheckResourceAttr(resName, "engine", strings.Split(engineVersion, "/")[0]),
					resource.TestCheckResourceAttr(resName, "version", strings.Split(engineVersion, "/")[1]),
				),
			},
		},
	})
}
