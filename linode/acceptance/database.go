package acceptance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func CheckMySQLDatabaseExists(name string, db *linodego.MySQLDatabase) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

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
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

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
