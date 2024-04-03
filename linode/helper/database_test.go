//go:build unit

package helper_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"reflect"
	"testing"
)

func TestResourceDatabaseMySQL_expandFlatten(t *testing.T) {
	data := linodego.MySQLDatabaseMaintenanceWindow{
		DayOfWeek: linodego.DatabaseMaintenanceDayWednesday,
		Duration:  1,
		Frequency: linodego.DatabaseMaintenanceFrequencyWeekly,
		HourOfDay: 5,
	}

	dataFlattened := helper.FlattenMaintenanceWindow(data)

	dataExpanded, err := helper.ExpandMaintenanceWindow(dataFlattened)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(dataExpanded, data) {
		t.Fatalf("maintenance window mismatch: %s", cmp.Diff(dataExpanded, data))
	}
}
