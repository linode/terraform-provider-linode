// //go:build unit

package instance

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func TestExpandInstanceConfigDeviceMap(t *testing.T) {
	deviceMapInput := map[string]any{
		"sda": []any{
			map[string]any{
				"disk_id":   124458,
				"volume_id": nil,
			},
		},
		"sdb": []any{
			map[string]any{
				"disk_id":   124459,
				"volume_id": nil,
			},
		},
	}

	diskIDLabelMap := map[string]int{
		"example_label_sda": 124458,
		"example_label_sdb": 124459,
	}

	deviceMap, err := expandInstanceConfigDeviceMap(deviceMapInput, diskIDLabelMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if deviceMap == nil {
		t.Error("Expected deviceMap to not be nil")
	}

	// Assert the DiskID for SDA and SDB
	expectedDiskID_SDA := 124458
	if deviceMap.SDA.DiskID != expectedDiskID_SDA {
		t.Errorf("Expected DiskID %d for SDA, but got %d", expectedDiskID_SDA, deviceMap.SDA.DiskID)
	}

	expectedDiskID_SDB := 124459
	if deviceMap.SDB.DiskID != expectedDiskID_SDB {
		t.Errorf("Expected DiskID %d for SDB, but got %d", expectedDiskID_SDB, deviceMap.SDB.DiskID)
	}

	// Assert that other device slots are nil
	if deviceMap.SDC != nil || deviceMap.SDD != nil || deviceMap.SDE != nil ||
		deviceMap.SDF != nil || deviceMap.SDG != nil || deviceMap.SDH != nil {
		t.Errorf("Expected all other device slots to be nil")
	}
}

func TestExpandInstanceConfigDevice(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		want *linodego.InstanceConfigDevice
	}{
		{
			name: "Valid DiskID",
			m: map[string]any{
				"disk_id": 123,
			},
			want: &linodego.InstanceConfigDevice{
				DiskID: 123,
			},
		},
		{
			name: "Valid VolumeID",
			m: map[string]any{
				"volume_id": 456,
			},
			want: &linodego.InstanceConfigDevice{
				VolumeID: 456,
			},
		},
		{
			name: "Invalid IDs",
			m: map[string]any{
				"disk_id":   0,
				"volume_id": -1,
			},
			want: nil,
		},
		{
			name: "No IDs",
			m:    map[string]any{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandInstanceConfigDevice(tt.m)
			if got != nil && tt.want == nil || got == nil && tt.want != nil {
				t.Errorf("expandInstanceConfigDevice() = %v, want %v", got, tt.want)
			} else if got != nil && tt.want != nil {
				if got.DiskID != tt.want.DiskID || got.VolumeID != tt.want.VolumeID {
					t.Errorf("expandInstanceConfigDevice() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestExpandInstanceACLPAlertsOpts(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want *linodego.InstanceACLPAlertsOptions
	}{
		{
			name: "Valid system_alerts and user_alerts",
			in: map[string]any{
				"system_alerts": schema.NewSet(schema.HashInt, []interface{}{1, 2}),
				"user_alerts":   schema.NewSet(schema.HashInt, []any{3, 4}),
			},
			want: func() *linodego.InstanceACLPAlertsOptions {
				return &linodego.InstanceACLPAlertsOptions{
					SystemAlerts: &[]int{1, 2},
					UserAlerts:   &[]int{3, 4},
				}
			}(),
		},
		{
			name: "Empty system_alerts and user_alerts",
			in: map[string]any{
				"system_alerts": schema.NewSet(schema.HashInt, []any{}),
				"user_alerts":   schema.NewSet(schema.HashInt, []any{}),
			},
			want: func() *linodego.InstanceACLPAlertsOptions {
				return &linodego.InstanceACLPAlertsOptions{
					SystemAlerts: &[]int{},
					UserAlerts:   &[]int{},
				}
			}(),
		},
		{
			name: "Absent system_alerts and user_alerts",
			in:   map[string]any{},
			want: &linodego.InstanceACLPAlertsOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandInstanceACLPAlertsOpts(tt.in)
			if got == nil {
				t.Fatalf("expected non-nil result")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("unexpected result. want: %v got: %v", tt.want, got)
			}
		})
	}
}

func TestExpandInstanceAlertsUpdateOpts(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want *linodego.InstanceAlert
	}{
		{
			name: "Valid legacy and ACLP alerts",
			in: map[string]any{
				"cpu":            90,
				"io":             1000,
				"network_in":     10,
				"network_out":    11,
				"transfer_quota": 80,
				"system_alerts": schema.NewSet(
					schema.HashInt,
					[]any{7, 8},
				),
				"user_alerts": schema.NewSet(schema.HashInt, []any{100}),
			},
			want: func() *linodego.InstanceAlert {
				return &linodego.InstanceAlert{
					CPU:           90,
					IO:            1000,
					NetworkIn:     10,
					NetworkOut:    11,
					TransferQuota: 80,
					SystemAlerts:  &[]int{7, 8},
					UserAlerts:    &[]int{100},
				}
			}(),
		},
		{
			name: "Only ACLP alerts provided",
			in: map[string]any{
				"system_alerts": schema.NewSet(schema.HashInt, []any{1}),
				"user_alerts":   schema.NewSet(schema.HashInt, []any{}),
			},
			want: func() *linodego.InstanceAlert {
				return &linodego.InstanceAlert{
					SystemAlerts: &[]int{1},
					UserAlerts:   &[]int{},
				}
			}(),
		},
		{
			name: "No alerts provided",
			in:   map[string]any{},
			want: &linodego.InstanceAlert{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandInstanceAlertsUpdateOpts(tt.in)
			if got == nil {
				t.Fatalf("expected non-nil result")
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("unexpected result. want: %v got: %v", tt.want, got)
			}
		})
	}
}
