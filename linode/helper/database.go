package helper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

var ValidDatabaseTypes = []string{"postgresql", "mysql"}

var UpdateObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"day_of_week":   types.StringType,
		"duration":      types.Int64Type,
		"frequency":     types.StringType,
		"hour_of_day":   types.Int64Type,
		"week_of_month": types.Int64Type,
	},
}

func ResolveValidDBEngine(
	ctx context.Context, client linodego.Client, engine string,
) (*linodego.DatabaseEngine, error) {
	filter := linodego.Filter{}
	filter.AddField(linodego.Eq, "engine", engine)

	filterBytes, err := filter.MarshalJSON()
	if err != nil {
		return nil, err
	}

	engines, err := client.ListDatabaseEngines(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return nil, err
	}

	if len(engines) < 1 {
		return nil, fmt.Errorf("no db engines were found")
	}

	return &engines[0], nil
}

var dayOfWeekStrToKey = map[string]linodego.DatabaseDayOfWeek{
	"sunday":    linodego.DatabaseMaintenanceDaySunday,
	"monday":    linodego.DatabaseMaintenanceDayMonday,
	"tuesday":   linodego.DatabaseMaintenanceDayTuesday,
	"wednesday": linodego.DatabaseMaintenanceDayWednesday,
	"thursday":  linodego.DatabaseMaintenanceDayThursday,
	"friday":    linodego.DatabaseMaintenanceDayFriday,
	"saturday":  linodego.DatabaseMaintenanceDaySaturday,
}

var dayOfWeekKeyToStr = map[linodego.DatabaseDayOfWeek]string{
	linodego.DatabaseMaintenanceDaySunday:    "sunday",
	linodego.DatabaseMaintenanceDayMonday:    "monday",
	linodego.DatabaseMaintenanceDayTuesday:   "tuesday",
	linodego.DatabaseMaintenanceDayWednesday: "wednesday",
	linodego.DatabaseMaintenanceDayThursday:  "thursday",
	linodego.DatabaseMaintenanceDayFriday:    "friday",
	linodego.DatabaseMaintenanceDaySaturday:  "saturday",
}

func ExpandDayOfWeek(day string) (linodego.DatabaseDayOfWeek, error) {
	result, ok := dayOfWeekStrToKey[day]
	if !ok {
		return 0, fmt.Errorf("invalid day of week: %s", day)
	}

	return result, nil
}

func FlattenDayOfWeek(day linodego.DatabaseDayOfWeek) string {
	return dayOfWeekKeyToStr[day]
}

func CreateLegacyDatabaseEngineSlug(engine, version string) string {
	return fmt.Sprintf("%s/%s", engine, version)
}

func CreateDatabaseEngineSlug(engine, version string) string {
	return fmt.Sprintf("%s/%s", engine, strings.Split(version, ".")[0])
}

func ParseDatabaseEngineSlug(engineID string) (string, string, error) {
	components := strings.Split(engineID, "/")
	if len(components) != 2 {
		return "", "", fmt.Errorf("invalid number of components: %d != 2", len(components))
	}

	return components[0], components[1], nil
}

func FlattenMaintenanceWindow(window linodego.MySQLDatabaseMaintenanceWindow) map[string]interface{} {
	result := make(map[string]interface{})

	result["day_of_week"] = FlattenDayOfWeek(window.DayOfWeek)
	result["duration"] = window.Duration
	result["frequency"] = string(window.Frequency)
	result["hour_of_day"] = window.HourOfDay

	// Nullable
	if window.WeekOfMonth != nil {
		result["week_of_month"] = window.WeekOfMonth
	}

	return result
}

func ExpandMaintenanceWindow(window map[string]interface{}) (linodego.DatabaseMaintenanceWindow, error) {
	result := linodego.DatabaseMaintenanceWindow{
		Duration:    window["duration"].(int),
		Frequency:   linodego.DatabaseMaintenanceFrequency(window["frequency"].(string)),
		HourOfDay:   window["hour_of_day"].(int),
		WeekOfMonth: nil,
	}

	dayOfWeek, err := ExpandDayOfWeek(window["day_of_week"].(string))
	if err != nil {
		return result, err
	}
	result.DayOfWeek = dayOfWeek

	if val, ok := window["week_of_month"]; ok && val.(int) > 0 {
		valInt := val.(int)
		result.WeekOfMonth = &valInt
	}

	return result, nil
}

func WaitForDatabaseUpdated(ctx context.Context, client linodego.Client, dbID int,
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

func FlattenDatabaseMaintenanceWindow(ctx context.Context, maintenance linodego.DatabaseMaintenanceWindow) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["day_of_week"] = types.StringValue(FlattenDayOfWeek(maintenance.DayOfWeek))
	result["duration"] = types.Int64Value(int64(maintenance.Duration))
	result["frequency"] = types.StringValue(string(maintenance.Frequency))
	result["hour_of_day"] = types.Int64Value(int64(maintenance.HourOfDay))
	result["week_of_month"] = IntPointerValueWithDefault(maintenance.WeekOfMonth)

	obj, diag := types.ObjectValue(UpdateObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		UpdateObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

// DatabaseStatusIsSuspended returns whether the given status is a "suspended" state.
func DatabaseStatusIsSuspended(databaseStatus linodego.DatabaseStatus) bool {
	return databaseStatus == linodego.DatabaseStatusSuspended || databaseStatus == linodego.DatabaseStatusSuspending
}

type databaseSuspensionFunc func(context.Context, int) error

// ReconcileDatabaseSuspensionSync synchronously suspends or resumes a database using the given functions
// depending on the given current and desired suspension statuses.
func ReconcileDatabaseSuspensionSync(
	ctx context.Context,
	client *linodego.Client,
	databaseID int,
	databaseEngine linodego.DatabaseEngineType,
	databaseSuspended bool,
	desiredSuspensionStatus bool,
	timeout time.Duration,
) error {
	var suspend, resume, targetOperation databaseSuspensionFunc
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
