package accountsettings

import (
	"github.com/linode/linodego"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseAccountSettings(t *testing.T) {
	// Create mock AccountSettings data
	mockEmail := "test@example.com"
	longviewSubscriptionValue := "longview-3"
	objectStorageValue := "active"
	backupsEnabledValue := true
	managedValue := true
	networkHelperValue := false

	mockSettings := &linodego.AccountSettings{
		BackupsEnabled:       backupsEnabledValue,
		Managed:              managedValue,
		NetworkHelper:        networkHelperValue,
		LongviewSubscription: &longviewSubscriptionValue,
		ObjectStorage:        &objectStorageValue,
	}

	// Create a mock AccountSettingsModel instance
	model := &AccountSettingsModel{}

	// Call the parseAccountSettings function
	model.parseAccountSettings(mockEmail, mockSettings)

	// Check if the fields in the model have been populated correctly
	if model.ID != types.StringValue(mockEmail) {
		t.Errorf("Expected ID to be %s, but got %s", mockEmail, model.ID)
	}

	if model.LongviewSubscription != types.StringValue("longview-3") {
		t.Errorf("Expected LongviewSubscription to be %s, but got %s", "longview-3", model.LongviewSubscription)
	}

	if model.ObjectStorage != types.StringValue("active") {
		t.Errorf("Expected ObjectStorage to be %s, but got %s", "active", model.ObjectStorage)
	}

	if model.BackupsEnabed != types.BoolValue(true) {
		t.Errorf("Expected BackupsEnabed to be %v, but got %v", true, model.BackupsEnabed)
	}

	if model.Managed != types.BoolValue(true) {
		t.Errorf("Expected Managed to be %v, but got %v", true, model.Managed)
	}

	if model.NetworkHelper != types.BoolValue(false) {
		t.Errorf("Expected NetworkHelper to be %v, but got %v", false, model.NetworkHelper)
	}
}
