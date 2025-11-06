package acceptance

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
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

var (
	DatabaseGetPathRegex           = regexp.MustCompile("/databases/\\S+/instances/[0-9d]+$")
	DatabasePendingUpdatesOverride = []map[string]any{
		{
			"description": "pending update 1",
		},
		{
			"deadline":    "2025-12-13T07:04:07",
			"description": "pending update 2 with deadline",
			"planned_for": "2025-12-13T07:04:07",
		},
		{
			"description": "pending update 3",
		},
	}
	DatabasePendingUpdatesSetExact = knownvalue.SetExact(
		[]knownvalue.Check{
			knownvalue.ObjectExact(
				map[string]knownvalue.Check{
					"deadline":    knownvalue.Null(),
					"description": knownvalue.StringExact(DatabasePendingUpdatesOverride[0]["description"].(string)),
					"planned_for": knownvalue.Null(),
				},
			),
			knownvalue.ObjectExact(
				map[string]knownvalue.Check{
					"deadline":    knownvalue.NotNull(),
					"description": knownvalue.StringExact(DatabasePendingUpdatesOverride[1]["description"].(string)),
					"planned_for": knownvalue.NotNull(),
				},
			),
			knownvalue.ObjectExact(
				map[string]knownvalue.Check{
					"deadline":    knownvalue.Null(),
					"description": knownvalue.StringExact(DatabasePendingUpdatesOverride[2]["description"].(string)),
					"planned_for": knownvalue.Null(),
				},
			),
		},
	)
)

func NewClientWithDatabasePendingUpdates(
	t *testing.T,
) *linodego.Client {
	return NewResponseOverrideClient(
		t,
		func(response *http.Response) bool {
			return response.Request.Method == "GET" && DatabaseGetPathRegex.MatchString(response.Request.RequestURI)
		},
		func(t *testing.T, responseBody map[string]any) {
			updates, ok := responseBody["updates"]
			if !ok {
				responseBody["updates"] = make(map[string]any)
				updates = responseBody["updates"]
			}

			updates.(map[string]any)["pending"] = DatabasePendingUpdatesOverride
		},
	)
}
