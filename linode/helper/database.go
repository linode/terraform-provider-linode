package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/linode/linodego"
)

var ValidDatabaseTypes = []string{"mongodb", "postgresql", "mysql"}

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

func CreateDatabaseEngineSlug(engine, version string) string {
	return fmt.Sprintf("%s/%s", engine, version)
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
	err = WaitForDatabaseStatusWithRetries(ctx, client, dbID, dbType,
		linodego.DatabaseStatusActive, timeoutSeconds)
	if err != nil {
		return fmt.Errorf("failed to wait for database active: %s", err)
	}

	return nil
}

// WaitForDatabaseStatusWithRetries waits for a database to reach a certain status
// with tolerance for intermittent API errors.
func WaitForDatabaseStatusWithRetries(
	ctx context.Context, client linodego.Client, dbID int, dbEngine linodego.DatabaseEngineType,
	status linodego.DatabaseStatus, timeoutSeconds int,
) error {
	return RunWithStatusRetries([]int{502}, 10, time.Second*5, func() error {
		return client.WaitForDatabaseStatus(ctx, dbID, dbEngine, status, timeoutSeconds)
	})
}
