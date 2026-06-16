//go:build unit

package monitoralertchannels

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func TestFlattenMonitorAlertChannel_AllFields(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	updated := created.Add(2 * time.Hour)

	input := &linodego.AlertChannel{
		ID:          123,
		Label:       "test-channel",
		Type:        linodego.UserAlertChannel,
		ChannelType: linodego.EmailAlertNotification,
		Created:     &created,
		Updated:     &updated,
		CreatedBy:   "creator",
		UpdatedBy:   "updater",
		Alerts: linodego.AlertsInfo{
			URL:        "/v4/monitor/alerts",
			Type:       "cpu",
			AlertCount: 5,
		},
		Details: linodego.ChannelDetails{
			Email: &linodego.EmailChannelDetails{
				Usernames:     []string{"alice", "bob"},
				RecipientType: "user",
			},
		},
	}

	got := flattenMonitorAlertChannel(context.Background(), input)

	require.Equal(t, int64(123), got.ID.ValueInt64())
	require.Equal(t, "test-channel", got.Label.ValueString())
	require.Equal(t, "user", got.Type.ValueString())
	require.Equal(t, "email", got.ChannelType.ValueString())
	require.Equal(t, "creator", got.CreatedBy.ValueString())
	require.Equal(t, "updater", got.UpdatedBy.ValueString())

	require.NotNil(t, got.Alerts)
	require.Equal(t, "/v4/monitor/alerts", got.Alerts.URL.ValueString())
	require.Equal(t, "cpu", got.Alerts.Type.ValueString())
	require.Equal(t, int64(5), got.Alerts.AlertCount.ValueInt64())

	require.NotNil(t, got.Details)
	require.NotNil(t, got.Details.Email)
	var usernames []string
	diags := got.Details.Email.Usernames.ElementsAs(context.Background(), &usernames, false)
	require.False(t, diags.HasError(), "failed to decode usernames")
	require.Equal(t, []string{"alice", "bob"}, usernames)
	require.Equal(t, "user", got.Details.Email.RecipientType.ValueString())
}

func TestFlattenMonitorAlertChannel_NilNested(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	input := &linodego.AlertChannel{
		ID:          1,
		Label:       "system-default",
		Type:        linodego.SystemAlertChannel,
		ChannelType: linodego.EmailAlertNotification,
		Created:     &now,
		Updated:     &now,
		CreatedBy:   "system",
		UpdatedBy:   "system",
		Alerts: linodego.AlertsInfo{
			URL:        "/v4/monitor/alerts",
			Type:       "disk",
			AlertCount: 0,
		},
	}

	got := flattenMonitorAlertChannel(context.Background(), input)
	require.Nil(t, got.Details)
}

func TestParseMonitorAlertChannels(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	model := MonitorAlertChannelFilterModel{}

	model.parseMonitorAlertChannels(context.Background(), []linodego.AlertChannel{
		{
			ID:          10,
			Label:       "one",
			Type:        linodego.UserAlertChannel,
			ChannelType: linodego.EmailAlertNotification,
			Created:     &now,
			Updated:     &now,
			CreatedBy:   "u1",
			UpdatedBy:   "u1",
			Alerts: linodego.AlertsInfo{
				URL:        "/alerts/1",
				Type:       "cpu",
				AlertCount: 1,
			},
		},
		{
			ID:          20,
			Label:       "two",
			Type:        linodego.SystemAlertChannel,
			ChannelType: linodego.EmailAlertNotification,
			Created:     &now,
			Updated:     &now,
			CreatedBy:   "system",
			UpdatedBy:   "system",
			Alerts: linodego.AlertsInfo{
				URL:        "/alerts/2",
				Type:       "memory",
				AlertCount: 2,
			},
		},
	})

	require.Len(t, model.MonitorAlertChannels, 2)
	require.Equal(t, int64(10), model.MonitorAlertChannels[0].ID.ValueInt64())
	require.Equal(t, int64(20), model.MonitorAlertChannels[1].ID.ValueInt64())
}
