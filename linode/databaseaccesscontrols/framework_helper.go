package databaseaccesscontrols

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func updateDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine string,
	id int,
	allowList []string,
) error {
	var createdDate *time.Time

	switch engine {
	case "mysql":
		db, err := client.UpdateMySQLDatabase(ctx, id, linodego.MySQLUpdateOptions{
			AllowList: &allowList,
		})
		if err != nil {
			return err
		}

		createdDate = db.Created
	case "postgresql":
		db, err := client.UpdatePostgresDatabase(ctx, id, linodego.PostgresUpdateOptions{
			AllowList: &allowList,
		})
		if err != nil {
			return err
		}

		createdDate = db.Created

	default:
		return fmt.Errorf("invalid database engine: %s", engine)
	}

	return helper.WaitForDatabaseUpdated(ctx, *client, id, linodego.DatabaseEngineType(engine),
		createdDate, 400)
}

func getDBAllowListByEngine(
	ctx context.Context,
	client *linodego.Client,
	engine string,
	id int,
) ([]string, error) {
	switch engine {
	case "mysql":
		db, err := client.GetMySQLDatabase(ctx, id)
		if err != nil {
			return nil, err
		}

		return db.AllowList, nil
	case "postgresql":
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
