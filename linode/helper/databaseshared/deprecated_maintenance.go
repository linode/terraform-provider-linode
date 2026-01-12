package databaseshared

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

var UpdateObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"day_of_week":   types.StringType,
		"duration":      types.Int64Type,
		"frequency":     types.StringType,
		"hour_of_day":   types.Int64Type,
		"week_of_month": types.Int64Type,
	},
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

func FlattenMaintenanceWindow(window linodego.MySQLDatabaseMaintenanceWindow) map[string]any {
	result := make(map[string]any)

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

func ExpandMaintenanceWindow(window map[string]any) (linodego.DatabaseMaintenanceWindow, error) {
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

func FlattenDatabaseMaintenanceWindow(ctx context.Context, maintenance linodego.DatabaseMaintenanceWindow) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["day_of_week"] = types.StringValue(FlattenDayOfWeek(maintenance.DayOfWeek))
	result["duration"] = types.Int64Value(int64(maintenance.Duration))
	result["frequency"] = types.StringValue(string(maintenance.Frequency))
	result["hour_of_day"] = types.Int64Value(int64(maintenance.HourOfDay))
	result["week_of_month"] = helper.IntPointerValueWithDefault(maintenance.WeekOfMonth)

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
