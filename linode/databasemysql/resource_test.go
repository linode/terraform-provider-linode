package databasemysql_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databasemysql/tmpl"
)

// TODO: resolve this dynamically
const engineVersion = "mysql/8.0.26"

func init() {
	resource.AddTestSweepers("linode_database_mysql", &resource.Sweeper{
		Name: "linode_database_mysql",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	_, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	_ = acceptance.SweeperListOptions(prefix, "database_mysql")

	// TODO: Sweep databases if acceptance.ShouldSweep(prefix, db.Label)

	return nil
}

func TestAccResourceDatabaseMySQL_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_database_mysql.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, dbName, engineVersion),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "engine", engineVersion),
					resource.TestCheckResourceAttr(resName, "label", dbName),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),

					resource.TestCheckResourceAttrSet(resName, "allow_list"),
					resource.TestCheckResourceAttr(resName, "cluster_size", "1"),
					resource.TestCheckResourceAttr(resName, "encrypted", "false"),
					resource.TestCheckResourceAttr(resName, "replication_type", "none"),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "false"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttrSet(resName, "host_secondary"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "version", strings.Split(engineVersion, "/")[0]),
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

func TestAccResourceDatabaseMySQL_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_database_mysql.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, tmpl.TemplateData{
					Engine:          engineVersion,
					Label:           dbName,
					AllowedIP:       "0.0.0.0/0",
					ClusterSize:     3,
					Encrypted:       true,
					ReplicationType: "asynch",
					SSLConnection:   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "engine", engineVersion),
					resource.TestCheckResourceAttr(resName, "label", dbName),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),

					resource.TestCheckResourceAttr(resName, "cluster_size", "3"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "replication_type", "asynch"),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttrSet(resName, "host_secondary"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "version", strings.Split(engineVersion, "/")[0]),
				),
			},
			{
				Config: tmpl.Complex(t, tmpl.TemplateData{
					Engine:          engineVersion,
					Label:           dbName + "updated",
					AllowedIP:       "192.0.2.1/32",
					ClusterSize:     3,
					Encrypted:       true,
					ReplicationType: "asynch",
					SSLConnection:   true,
				}),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "engine", engineVersion),
					resource.TestCheckResourceAttr(resName, "label", dbName+"updated"),
					resource.TestCheckResourceAttr(resName, "region", "us-southeast"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),

					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "192.0.2.1/32"),

					resource.TestCheckResourceAttr(resName, "cluster_size", "3"),
					resource.TestCheckResourceAttr(resName, "encrypted", "true"),
					resource.TestCheckResourceAttr(resName, "replication_type", "asynch"),
					resource.TestCheckResourceAttr(resName, "ssl_connection", "true"),

					resource.TestCheckResourceAttrSet(resName, "ca_cert"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "host_primary"),
					resource.TestCheckResourceAttrSet(resName, "host_secondary"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "status", "active"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttrSet(resName, "root_password"),
					resource.TestCheckResourceAttr(resName, "version", strings.Split(engineVersion, "/")[0]),
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

func checkMySQLDatabaseExists(s *terraform.State) error {
	// TODO

	return nil
}

func checkDestroy(s *terraform.State) error {
	// TODO
	return nil
}
