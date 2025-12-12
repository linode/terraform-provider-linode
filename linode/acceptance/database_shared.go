package acceptance

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/linode/linodego"
)

var (
	DatabaseGetPathRegex           = regexp.MustCompile(`/databases/\S+/instances/\d+$`)
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
		// This duplicate will be filtered out
		{
			"description": "pending update 1",
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

// NewClientWithDatabasePendingUpdates returns a Linode client with injected pending updates
// for managed databases.
func NewClientWithDatabasePendingUpdates(
	t *testing.T,
) *linodego.Client {
	return NewResponseOverrideClient(
		t,
		func(response *http.Response) bool {
			return response.Request.Method == "GET" && DatabaseGetPathRegex.MatchString(response.Request.URL.Path)
		},
		func(responseBody map[string]any) error {
			updates, ok := responseBody["updates"]
			if !ok {
				responseBody["updates"] = make(map[string]any)
				updates = responseBody["updates"]
			}

			updates.(map[string]any)["pending"] = DatabasePendingUpdatesOverride

			return nil
		},
	)
}
