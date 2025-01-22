//go:build integration || databasemysqlv2

package databasemysqlv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databasemysqlv2/tmpl"
)

func TestAccDataSource_basic(t *testing.T) {
	t.Parallel()

	label := acctest.RandomWithPrefix("tf_test")
	dataSourceName := "data.linode_database_mysql_v2.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Data(t, tmpl.TemplateData{
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
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),

					resource.TestCheckResourceAttrSet(dataSourceName, "ca_cert"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_size", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "created"),
					resource.TestCheckResourceAttr(dataSourceName, "encrypted", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "engine", "mysql"),
					resource.TestCheckResourceAttr(dataSourceName, "engine_id", testEngine),
					resource.TestCheckNoResourceAttr(dataSourceName, "fork_restore_time"),
					resource.TestCheckNoResourceAttr(dataSourceName, "fork_source"),
					resource.TestCheckResourceAttrSet(dataSourceName, "host_primary"),
					resource.TestCheckResourceAttr(dataSourceName, "label", label),
					resource.TestCheckResourceAttrSet(dataSourceName, "members.%"),
					resource.TestCheckResourceAttrSet(dataSourceName, "root_password"),
					resource.TestCheckResourceAttrSet(dataSourceName, "root_username"),
					resource.TestCheckResourceAttr(dataSourceName, "platform", "rdbms-default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "port"),
					resource.TestCheckResourceAttr(dataSourceName, "region", testRegion),
					resource.TestCheckResourceAttr(dataSourceName, "ssl_connection", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "active"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),

					resource.TestCheckResourceAttr(dataSourceName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "allow_list.0", "10.0.0.4/32"),

					resource.TestCheckResourceAttr(dataSourceName, "updates.hour_of_day", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.day_of_week", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.duration", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "updates.frequency", "weekly"),

					resource.TestCheckResourceAttr(dataSourceName, "pending_updates.#", "0"),
				),
			},
		},
	})
}
