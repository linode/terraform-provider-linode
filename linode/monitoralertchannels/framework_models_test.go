//go:build unit

package monitoralertchannels

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego"
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
		Content: linodego.ChannelContent{
			Email: &linodego.EmailChannelContent{
				EmailAddresses: []string{"a@example.com", "b@example.com"},
			},
		},
		Details: linodego.ChannelDetails{
			Email: &linodego.EmailChannelDetails{
				Usernames:     []string{"alice", "bob"},
				RecipientType: "user",
			},
		},
	}

	got := flattenMonitorAlertChannel(context.Background(), input)

	if got.ID.ValueInt64() != 123 {
		t.Fatalf("expected id=123, got %d", got.ID.ValueInt64())
	}
	if got.Label.ValueString() != "test-channel" {
		t.Fatalf("expected label=test-channel, got %s", got.Label.ValueString())
	}
	if got.Type.ValueString() != "user" {
		t.Fatalf("expected type=user, got %s", got.Type.ValueString())
	}
	if got.ChannelType.ValueString() != "email" {
		t.Fatalf("expected channel_type=email, got %s", got.ChannelType.ValueString())
	}
	if got.CreatedBy.ValueString() != "creator" {
		t.Fatalf("expected created_by=creator, got %s", got.CreatedBy.ValueString())
	}
	if got.UpdatedBy.ValueString() != "updater" {
		t.Fatalf("expected updated_by=updater, got %s", got.UpdatedBy.ValueString())
	}

	if got.Alerts == nil {
		t.Fatal("expected alerts to be set")
	}
	if got.Alerts.URL.ValueString() != "/v4/monitor/alerts" {
		t.Fatalf("expected alerts.url=/v4/monitor/alerts, got %s", got.Alerts.URL.ValueString())
	}
	if got.Alerts.Type.ValueString() != "cpu" {
		t.Fatalf("expected alerts.type=cpu, got %s", got.Alerts.Type.ValueString())
	}
	if got.Alerts.AlertCount.ValueInt64() != 5 {
		t.Fatalf("expected alerts.alert_count=5, got %d", got.Alerts.AlertCount.ValueInt64())
	}

	if got.Content == nil || got.Content.Email == nil {
		t.Fatal("expected content.email to be set")
	}
	var emails []string
	diags := got.Content.Email.EmailAddresses.ElementsAs(context.Background(), &emails, false)
	if diags.HasError() {
		t.Fatalf("failed to decode email_addresses: %v", diags)
	}
	if len(emails) != 2 || emails[0] != "a@example.com" || emails[1] != "b@example.com" {
		t.Fatalf("unexpected email_addresses: %#v", emails)
	}

	if got.Details == nil || got.Details.Email == nil {
		t.Fatal("expected details.email to be set")
	}
	var usernames []string
	diags = got.Details.Email.Usernames.ElementsAs(context.Background(), &usernames, false)
	if diags.HasError() {
		t.Fatalf("failed to decode usernames: %v", diags)
	}
	if len(usernames) != 2 || usernames[0] != "alice" || usernames[1] != "bob" {
		t.Fatalf("unexpected usernames: %#v", usernames)
	}
	if got.Details.Email.RecipientType.ValueString() != "user" {
		t.Fatalf("expected recipient_type=user, got %s", got.Details.Email.RecipientType.ValueString())
	}
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
	if got.Content != nil {
		t.Fatalf("expected content=nil, got %#v", got.Content)
	}
	if got.Details != nil {
		t.Fatalf("expected details=nil, got %#v", got.Details)
	}
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

	if len(model.MonitorAlertChannels) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(model.MonitorAlertChannels))
	}
	if model.MonitorAlertChannels[0].ID.ValueInt64() != 10 {
		t.Fatalf("expected first id=10, got %d", model.MonitorAlertChannels[0].ID.ValueInt64())
	}
	if model.MonitorAlertChannels[1].ID.ValueInt64() != 20 {
		t.Fatalf("expected second id=20, got %d", model.MonitorAlertChannels[1].ID.ValueInt64())
	}
}
