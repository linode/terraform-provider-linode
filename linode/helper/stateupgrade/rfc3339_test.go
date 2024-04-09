//go:build unit

package stateupgrade_test

import (
	"testing"
	"time"

	"github.com/linode/terraform-provider-linode/v2/linode/helper/stateupgrade"
)

func TestStackScriptUpgradeTimeFormat(t *testing.T) {
	t.Parallel()

	trueTime := time.Now().Round(time.Second)
	testTimeUpgrader(t, trueTime, trueTime.String())
	testTimeUpgrader(t, trueTime, trueTime.Format(time.RFC3339))
}

func testTimeUpgrader(t *testing.T, trueTime time.Time, timeString string) {
	t.Helper()

	upgradedRFC3339Time, err := stateupgrade.UpgradeTimeFormatToRFC3339(timeString)
	if err != nil {
		t.Fatal(err)
	}

	upgradedTime, diags := upgradedRFC3339Time.ValueRFC3339Time()
	if diags.HasError() {
		t.Fatal("failed to get ValueRFC3339Time")
	}

	if !upgradedTime.Equal(trueTime) {
		t.Fatal("time value not matched")
	}
}
