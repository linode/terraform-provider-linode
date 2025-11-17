package databaseshared

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// StatusIsSuspended returns whether the given status is a "suspended" state.
func StatusIsSuspended(databaseStatus linodego.DatabaseStatus) bool {
	return databaseStatus == linodego.DatabaseStatusSuspended || databaseStatus == linodego.DatabaseStatusSuspending
}

type suspensionFunc func(context.Context, int) error

// ReconcileSuspensionSync synchronously suspends or resumes a database using the given functions
// depending on the given current and desired suspension statuses.
func ReconcileSuspensionSync(
	ctx context.Context,
	client *linodego.Client,
	databaseID int,
	databaseEngine linodego.DatabaseEngineType,
	databaseSuspended bool,
	desiredSuspensionStatus bool,
	timeout time.Duration,
) error {
	var suspend, resume, targetOperation suspensionFunc
	var desiredStatus linodego.DatabaseStatus

	switch databaseEngine {
	case linodego.DatabaseEngineTypeMySQL:
		suspend, resume = client.SuspendMySQLDatabase, client.ResumeMySQLDatabase
	case linodego.DatabaseEngineTypePostgres:
		suspend, resume = client.SuspendPostgresDatabase, client.ResumePostgresDatabase
	}

	if databaseSuspended && !desiredSuspensionStatus {
		targetOperation = resume
		desiredStatus = linodego.DatabaseStatusActive
	} else if !databaseSuspended && desiredSuspensionStatus {
		targetOperation = suspend
		desiredStatus = linodego.DatabaseStatusSuspended
	}

	if targetOperation == nil {
		// Nothing to do here
		return nil
	}

	tflog.Debug(ctx, "Calling target function to reconcile database suspension")
	if err := targetOperation(ctx, databaseID); err != nil {
		return fmt.Errorf("failed to reconcile suspension of database: %w", err)
	}

	tflog.Debug(ctx, "client.WaitForDatabaseStatus(...)", map[string]any{
		"status": desiredStatus,
	})
	if err := client.WaitForDatabaseStatus(
		ctx,
		databaseID,
		databaseEngine,
		desiredStatus,
		int(timeout.Seconds()),
	); err != nil {
		return fmt.Errorf("failed to wait for database status %s: %w", desiredStatus, err)
	}

	return nil
}
