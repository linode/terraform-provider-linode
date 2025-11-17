package databaseshared

import (
	"context"
	"fmt"
	"time"

	"github.com/linode/linodego"
)

var ValidDatabaseTypes = []string{"postgresql", "mysql"}

func WaitForUpdated(ctx context.Context, client linodego.Client, dbID int,
	dbType linodego.DatabaseEngineType, minStart *time.Time, timeoutSeconds int,
) error {
	if minStart == nil {
		return fmt.Errorf("nil minimum starting time")
	}

	_, err := client.WaitForEventFinished(ctx, dbID, linodego.EntityDatabase,
		linodego.ActionDatabaseUpdate, *minStart, timeoutSeconds)
	if err != nil {
		return fmt.Errorf("failed to wait for database update: %s", err)
	}

	// Sometimes the event has finished but the status hasn't caught up
	err = client.WaitForDatabaseStatus(ctx, dbID, dbType,
		linodego.DatabaseStatusActive, timeoutSeconds)
	if err != nil {
		return fmt.Errorf("failed to wait for database active: %s", err)
	}

	return nil
}
