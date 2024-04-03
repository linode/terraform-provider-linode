//go:build integration || databaseaccesscontrols

package databaseaccesscontrols_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/databaseaccesscontrols/tmpl"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var (
	mysqlEngineVersion    string
	postgresEngineVersion string
	testRegion            string
)

func init() {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	mysqlEngineVersion = v.ID

	v, err = helper.ResolveValidDBEngine(context.Background(), *client, "postgresql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	postgresEngineVersion = v.ID

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Managed Databases"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceDatabaseAccessControls_MySQL(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	resName := "linode_database_access_controls.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.MySQL(t, dbName, mysqlEngineVersion, "0.0.0.0/0", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),
				),
			},
			{
				Config: tmpl.MySQL(t, dbName, mysqlEngineVersion, "192.168.0.25/32", testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "192.168.0.25/32"),
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

func TestAccResourceDatabaseAccessControls_PostgreSQL(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	resName := "linode_database_access_controls.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.PostgreSQL(t, dbName, postgresEngineVersion, "0.0.0.0/0", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),
				),
			},
			{
				Config: tmpl.PostgreSQL(t, dbName, postgresEngineVersion, "192.168.0.25/32", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "192.168.0.25/32"),
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
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetMySQLDatabase(context.Background(), id)
		if err != nil {
			return fmt.Errorf("error retrieving state of mysql database %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetMySQLDatabase(context.Background(), id)

		if err == nil {
			return fmt.Errorf("mysql database with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting mysql database with id %d", id)
		}
	}

	return nil
}
