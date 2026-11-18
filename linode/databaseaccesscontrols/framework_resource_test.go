//go:build integration || databaseaccesscontrols || dbaas_tests

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
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/databaseaccesscontrols/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
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

	v, err := databaseshared.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	mysqlEngineVersion = v.ID

	v, err = databaseshared.ResolveValidDBEngine(context.Background(), *client, "postgresql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	postgresEngineVersion = v.ID

	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityDBAAS}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceDatabaseAccessControls_MySQL(t *testing.T) {
	t.Parallel()

	resName := "linode_database_access_controls.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckMySQLDatabaseV2Destroy,
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
	t.Parallel()

	resName := "linode_database_access_controls.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckPostgreSQLDatabaseV2Destroy,
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

func TestAccResourceDatabaseAccessControls_Valkey(t *testing.T) {
	t.Parallel()

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}
	engine, err := databaseshared.ResolveValidDBEngine(context.Background(), *client, string(linodego.DatabaseEngineTypeValkey))
	if err != nil {
		t.Fatalf("failed to get Valkey database engine version: %v", err)
	}

	const resourceName = "linode_database_access_controls.foobar"
	databaseLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckValkeyDatabaseV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Valkey(t, databaseLabel, engine.ID, "192.0.2.10/32", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_type", "valkey"),
					resource.TestCheckResourceAttr(resourceName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_list.0", "192.0.2.10/32"),
				),
			},
			{
				Config: tmpl.Valkey(t, databaseLabel, engine.ID, "192.168.0.25/32", testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_type", "valkey"),
					resource.TestCheckResourceAttr(resourceName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "allow_list.0", "192.168.0.25/32"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkMySQLDatabaseExists(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql_v2" {
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
