package databaseaccesscontrols

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func updateDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine linodego.DatabaseEngineType,
	id int,
	allowList []string,
	timeoutSeconds int,
) error {
	currentAllowList, err := getDBAllowListByEngine(ctx, client, engine, id)
	if err != nil {
		return err
	}

	// Sort because order isn't important
	sort.Strings(allowList)
	sort.Strings(currentAllowList)

	// Nothing to do here
	if reflect.DeepEqual(allowList, currentAllowList) {
		return nil
	}

	updateDatabase := func() error {
		// Reuse the error handler for these functions
		switch engine {
		case linodego.DatabaseEngineTypeMySQL:
			_, err = client.UpdateMySQLDatabase(ctx, id, linodego.MySQLUpdateOptions{
				AllowList: &allowList,
			})
		case linodego.DatabaseEngineTypePostgres:
			_, err = client.UpdatePostgresDatabase(ctx, id, linodego.PostgresUpdateOptions{
				AllowList: &allowList,
			})
		default:
			return fmt.Errorf("invalid database engine: %s", engine)
		}

		if err != nil {
			return fmt.Errorf("failed to update allow_list for database %d: %w", id, err)
		}

		return nil
	}

	return helper.WaitForDatabaseUpdated(
		ctx,
		client,
		id,
		engine,
		timeoutSeconds,
		updateDatabase,
	)
}

func getDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine linodego.DatabaseEngineType,
	id int,
) ([]string, error) {
	switch engine {
	case linodego.DatabaseEngineTypeMySQL:
		db, err := client.GetMySQLDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	case linodego.DatabaseEngineTypePostgres:
		db, err := client.GetPostgresDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	}

	return nil, fmt.Errorf("invalid database type: %s", engine)
}

func formatID(dbID int, dbType linodego.DatabaseEngineType) string {
	return fmt.Sprintf("%d:%s", dbID, dbType)
}

func parseID(id string) (int, linodego.DatabaseEngineType, error) {
	split := strings.Split(id, ":")
	if len(split) != 2 {
		return 0, "", fmt.Errorf("invalid number of segments")
	}

	dbID, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, "", err
	}

	return dbID, linodego.DatabaseEngineType(split[1]), nil
}
