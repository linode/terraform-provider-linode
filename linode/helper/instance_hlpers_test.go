//go:build unit

package helper

import (
	"testing"

	"github.com/linode/linodego"
)

func TestExpandConfigInterface(t *testing.T) {
	interfaceInput := map[string]interface{}{
		"label":        "eth0.100",
		"purpose":      "vlan",
		"ipam_address": "192.168.1.2/24",
		"primary":      false,
	}

	interfaceResult := ExpandConfigInterface(interfaceInput)

	expectedLabel := "eth0.100"
	if interfaceResult.Label != expectedLabel {
		t.Errorf("Expected label %s, but got %s", expectedLabel, interfaceResult.Label)
	}

	expectedPurpose := linodego.InterfacePurposeVLAN
	if interfaceResult.Purpose != expectedPurpose {
		t.Errorf("Expected purpose %s, but got %s", expectedPurpose, interfaceResult.Purpose)
	}

	expectedIPAMAddress := "192.168.1.2/24"
	if interfaceResult.IPAMAddress != expectedIPAMAddress {
		t.Errorf("Expected IPAMAddress %s, but got %s", expectedIPAMAddress, interfaceResult.IPAMAddress)
	}
}

func TestFlattenConfigInterface(t *testing.T) {
	configInterface := linodego.InstanceConfigInterface{
		IPAMAddress: "192.168.1.1",
		Label:       "test-vlan",
		Purpose:     linodego.InterfacePurposeVLAN,
	}

	result := FlattenInterface(configInterface)

	expected := map[string]interface{}{
		"ipam_address": "192.168.1.1",
		"label":        "test-vlan",
		"purpose":      linodego.InterfacePurposeVLAN,
	}

	for key, expectedValue := range expected {
		if resultValue, ok := result[key]; !ok || resultValue != expectedValue {
			t.Errorf("Mismatch for key %s: Expected %v, but got %v", key, expectedValue, resultValue)
		}
	}
}
