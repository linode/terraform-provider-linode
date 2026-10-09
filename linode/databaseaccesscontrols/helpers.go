package databaseaccesscontrols

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-set/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
)

// updateDBAllowListByEngine updates the allow list of the given database,
// waiting for the resulting update event to finish.
func updateDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine string,
	id int,
	allowList []string,
) error {
	// Rather than using the state, we retrieve and compare DB allow lists here
	// to account for a case where the state is not populated during creation.
	//
	// This is necessary because database_update events are only created when the value
	// of allow_list is actually changed.
	oldAllowList, err := getDBAllowListByEngine(ctx, client, engine, id)
	if err != nil {
		return fmt.Errorf("failed to get allow_list for database: %w", err)
	}

	if set.From(oldAllowList).Equal(set.From(allowList)) {
		// Nothing to do here
		return nil
	}

	updatePoller, err := client.NewEventPoller(ctx, id, linodego.EntityDatabase, linodego.ActionDatabaseUpdate)
	if err != nil {
		return fmt.Errorf("failed to create update EventPoller: %w", err)
	}

	if allowList == nil {
		allowList = []string{}
	}

	switch engine {
	case "mysql":
		tflog.Debug(ctx, "client.UpdateMySQLDatabase(...)", map[string]any{
			"allow_list": allowList,
		})

		if _, err := client.UpdateMySQLDatabase(ctx, id, linodego.MySQLUpdateOptions{
			AllowList: allowList,
		}); err != nil {
			return err
		}

	case "postgresql":
		tflog.Debug(ctx, "client.UpdatePostgresDatabase(...)", map[string]any{
			"allow_list": allowList,
		})

		if _, err := client.UpdatePostgresDatabase(ctx, id, linodego.PostgresUpdateOptions{
			AllowList: allowList,
		}); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid database engine: %s", engine)
	}

	if _, err := updatePoller.WaitForFinished(ctx); err != nil {
		return fmt.Errorf("failed to wait for update event completion: %w", err)
	}

	return nil
}

// getDBAllowListByEngine returns the current allow list of the given database.
func getDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine string,
	id int,
) ([]string, error) {
	switch engine {
	case "mysql":
		tflog.Trace(ctx, "client.GetMySQLDatabase(...)")

		db, err := client.GetMySQLDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	case "postgresql":
		tflog.Trace(ctx, "client.GetPostgresDatabase(...)")

		db, err := client.GetPostgresDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	}

	return nil, fmt.Errorf("invalid database type: %s", engine)
}

func formatID(dbID int, dbType string) string {
	return fmt.Sprintf("%d:%s", dbID, dbType)
}

func parseID(id string) (int, string, error) {
	split := strings.Split(id, ":")
	if len(split) != 2 {
		return 0, "", fmt.Errorf("invalid number of segments")
	}

	dbID, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, "", err
	}

	return dbID, split[1], nil
}
