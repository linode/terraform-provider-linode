package acceptance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func CheckMySQLDatabaseExists(name string, db *linodego.MySQLDatabase) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetMySQLDatabase(context.Background(), id)
		if err != nil {
			return fmt.Errorf("error retrieving state of mysql database %s: %s", rs.Primary.Attributes["label"], err)
		}

		if db != nil {
			*db = *found
		}

		return nil
	}
}

func CheckPostgresDatabaseExists(name string, db *linodego.PostgresDatabase) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetPostgresDatabase(context.Background(), id)
		if err != nil {
			return fmt.Errorf("error retrieving state of postgres database %s: %s", rs.Primary.Attributes["label"], err)
		}

		if db != nil {
			*db = *found
		}

		return nil
	}
}

func CheckMySQLDatabaseV2Destroy(s *terraform.State) error {
	client := TestAccFrameworkProvider.Meta.Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql_v2" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse %v as int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("should not have Linode ID 0")
		}

		_, err = client.GetMySQLDatabase(context.Background(), id)

		if err == nil {
			return fmt.Errorf("should not find database ID %d existing after delete", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && !linodego.IsNotFound(apiErr) {
			return fmt.Errorf("failed to get database ID %d: %s", id, err)
		}
	}

	return nil
}

func CheckPostgreSQLDatabaseV2Destroy(s *terraform.State) error {
	client := TestAccFrameworkProvider.Meta.Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_postgresql_v2" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse %v as int", rs.Primary.ID)
		}

		if id == 0 {
			return fmt.Errorf("should not have Linode ID 0")
		}

		_, err = client.GetPostgresDatabase(context.Background(), id)

		if err == nil {
			return fmt.Errorf("should not find database ID %d existing after delete", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && !linodego.IsNotFound(apiErr) {
			return fmt.Errorf("failed to get database ID %d: %s", id, err)
		}
	}

	return nil
}
