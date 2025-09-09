//go:build unit

package maintenancepolicies

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseMaintenancePolicies(t *testing.T) {
	mockPolicies := []linodego.MaintenancePolicy{
		{
			Slug:                  "linode/migrate",
			Label:                 "Migrate",
			Description:           "Migrates the Linode to a new host while it remains fully operational. Recommended for maximizing availability.",
			Type:                  "migrate",
			NotificationPeriodSec: 3600,
			IsDefault:             true,
		},
		{
			Slug:                  "linode/power_off_on",
			Label:                 "Power Off/Power On",
			Description:           "Powers off the Linode at the start of the maintenance event and reboots it once the maintenance finishes. Recommended for maximizing performance.",
			Type:                  "power_off_on",
			NotificationPeriodSec: 1800,
			IsDefault:             false,
		},
		{
			Slug:                  "private/12345",
			Label:                 "Critical Workload - Avoid Migration",
			Description:           "Custom policy designed to power off and perform maintenance during user-defined windows only.",
			Type:                  "power_off_on",
			NotificationPeriodSec: 7200,
			IsDefault:             false,
		},
	}

	data := &MaintenancePolicyFilterModel{}
	diags := data.parseMaintenancePolicies(mockPolicies)
	assert.False(t, diags.HasError(), "Unexpected diagnostics error")

	assert.Equal(t, types.StringValue("linode/migrate"), data.MaintenancePolicies[0].Slug)
	assert.Equal(t, types.StringValue("Migrate"), data.MaintenancePolicies[0].Label)
	assert.Equal(
		t,
		types.StringValue("Migrates the Linode to a new host while it remains fully operational. Recommended for maximizing availability."),
		data.MaintenancePolicies[0].Description,
	)
	assert.Equal(t, types.StringValue("migrate"), data.MaintenancePolicies[0].Type)
	assert.Equal(t, types.Int64Value(3600), data.MaintenancePolicies[0].NotificationPeriodSec)
	assert.Equal(t, types.BoolValue(true), data.MaintenancePolicies[0].IsDefault)

	assert.Equal(t, types.StringValue("linode/power_off_on"), data.MaintenancePolicies[1].Slug)
	assert.Equal(t, types.StringValue("Power Off/Power On"), data.MaintenancePolicies[1].Label)
	assert.Equal(
		t,
		types.StringValue(
			"Powers off the Linode at the start of the maintenance event and reboots it once the maintenance finishes. Recommended for maximizing performance.",
		),
		data.MaintenancePolicies[1].Description,
	)
	assert.Equal(t, types.StringValue("power_off_on"), data.MaintenancePolicies[1].Type)
	assert.Equal(t, types.Int64Value(1800), data.MaintenancePolicies[1].NotificationPeriodSec)
	assert.Equal(t, types.BoolValue(false), data.MaintenancePolicies[1].IsDefault)

	assert.Equal(t, types.StringValue("private/12345"), data.MaintenancePolicies[2].Slug)
	assert.Equal(t, types.StringValue("Critical Workload - Avoid Migration"), data.MaintenancePolicies[2].Label)
	assert.Equal(
		t,
		types.StringValue("Custom policy designed to power off and perform maintenance during user-defined windows only."),
		data.MaintenancePolicies[2].Description,
	)
	assert.Equal(t, types.StringValue("power_off_on"), data.MaintenancePolicies[2].Type)
	assert.Equal(t, types.Int64Value(7200), data.MaintenancePolicies[2].NotificationPeriodSec)
	assert.Equal(t, types.BoolValue(false), data.MaintenancePolicies[2].IsDefault)
}
