package databasepostgresqlv2

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type ModelHosts struct {
	Primary   types.String `tfsdk:"primary"`
	Secondary types.String `tfsdk:"secondary"`
}

type ModelUpdates struct {
	DayOfWeek types.Int64  `tfsdk:"day_of_week"`
	Duration  types.Int64  `tfsdk:"duration"`
	Frequency types.String `tfsdk:"frequency"`
	HourOfDay types.Int64  `tfsdk:"hour_of_day"`
}

func (m ModelUpdates) ToLinodego(d diag.Diagnostics) *linodego.DatabaseMaintenanceWindow {
	return &linodego.DatabaseMaintenanceWindow{
		DayOfWeek: linodego.DatabaseDayOfWeek(helper.FrameworkSafeInt64ToInt(m.DayOfWeek.ValueInt64(), &d)),
		Duration:  helper.FrameworkSafeInt64ToInt(m.Duration.ValueInt64(), &d),
		Frequency: linodego.DatabaseMaintenanceFrequency(m.Frequency.ValueString()),
		HourOfDay: helper.FrameworkSafeInt64ToInt(m.HourOfDay.ValueInt64(), &d),
	}
}

type ModelPendingUpdate struct {
	Deadline    timetypes.RFC3339 `tfsdk:"deadline"`
	Description types.String      `tfsdk:"description"`
	PlannedFor  timetypes.RFC3339 `tfsdk:"planned_for"`
}

type Model struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`

	ID types.String `tfsdk:"id"`

	AllowList     types.Set         `tfsdk:"allow_list"`
	CACert        types.String      `tfsdk:"ca_cert"`
	ClusterSize   types.Int64       `tfsdk:"cluster_size"`
	Created       timetypes.RFC3339 `tfsdk:"created"`
	Encrypted     types.Bool        `tfsdk:"encrypted"`
	Engine        types.String      `tfsdk:"engine"`
	EngineID      types.String      `tfsdk:"engine_id"`
	HostPrimary   types.String      `tfsdk:"host_primary"`
	HostSecondary types.String      `tfsdk:"host_secondary"`
	Label         types.String      `tfsdk:"label"`
	Members       types.Map         `tfsdk:"members"`
	Platform      types.String      `tfsdk:"platform"`
	Port          types.Int64       `tfsdk:"port"`
	Region        types.String      `tfsdk:"region"`
	RootPassword  types.String      `tfsdk:"root_password"`
	RootUsername  types.String      `tfsdk:"root_username"`
	SSLConnection types.Bool        `tfsdk:"ssl_connection"`
	Status        types.String      `tfsdk:"status"`
	Type          types.String      `tfsdk:"type"`
	Updated       timetypes.RFC3339 `tfsdk:"updated"`
	Version       types.String      `tfsdk:"version"`

	// Fork-specific fields
	OldestRestoreTime timetypes.RFC3339 `tfsdk:"oldest_restore_time"`
	ForkSource        types.Int64       `tfsdk:"fork_source"`
	ForkRestoreTime   timetypes.RFC3339 `tfsdk:"fork_restore_time"`

	Updates        types.Object `tfsdk:"updates"`
	PendingUpdates types.Set    `tfsdk:"pending_updates"`
}

func (m *Model) Refresh(
	ctx context.Context,
	client *linodego.Client,
	dbID int,
	preserveKnown bool,
) (d diag.Diagnostics) {
	tflog.SetField(ctx, "id", dbID)

	tflog.Debug(ctx, "Refreshing the PostgreSQL database...")

	tflog.Debug(ctx, "client.GetPostgresDatabase(...)")
	db, err := client.GetPostgresDatabase(ctx, dbID)
	if err != nil {
		d.AddError("Failed to refresh PostgreSQL database", err.Error())
		return
	}

	tflog.Debug(ctx, "client.GetPostgresDatabaseSSL(...)")
	dbSSL, err := client.GetPostgresDatabaseSSL(ctx, dbID)
	if err != nil {
		d.AddError("Failed to refresh PostgreSQL database SSL", err.Error())
		return
	}

	tflog.Debug(ctx, "client.GetPostgresDatabaseCredentials(...)")
	dbCreds, err := client.GetPostgresDatabaseCredentials(ctx, dbID)
	if err != nil {
		d.AddError("Failed to refresh PostgreSQL database credentials", err.Error())
		return
	}

	m.Flatten(ctx, db, dbSSL, dbCreds, preserveKnown)
	return
}

func (m *Model) Flatten(
	ctx context.Context,
	db *linodego.PostgresDatabase,
	ssl *linodego.PostgresDatabaseSSL,
	creds *linodego.PostgresDatabaseCredential,
	preserveKnown bool,
) (d diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(db.ID), preserveKnown)

	m.CACert = helper.KeepOrUpdateString(m.CACert, string(ssl.CACertificate), preserveKnown)
	m.ClusterSize = helper.KeepOrUpdateInt64(m.ClusterSize, int64(db.ClusterSize), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(db.Created), preserveKnown)
	m.Encrypted = helper.KeepOrUpdateBool(m.Encrypted, db.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateString(m.Engine, db.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateString(
		m.EngineID,
		helper.CreateDatabaseEngineSlug(db.Engine, db.Version),
		preserveKnown,
	)
	m.HostPrimary = helper.KeepOrUpdateString(m.HostPrimary, db.Hosts.Primary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateString(m.HostSecondary, db.Hosts.Secondary, preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, db.Label, preserveKnown)
	m.OldestRestoreTime = helper.KeepOrUpdateValue(m.OldestRestoreTime, timetypes.NewRFC3339TimePointerValue(db.OldestRestoreTime), preserveKnown)
	m.Platform = helper.KeepOrUpdateString(m.Platform, string(db.Platform), preserveKnown)
	m.Port = helper.KeepOrUpdateInt64(m.Port, int64(db.Port), preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, db.Region, preserveKnown)
	m.RootPassword = helper.KeepOrUpdateString(m.RootPassword, creds.Password, preserveKnown)
	m.RootUsername = helper.KeepOrUpdateString(m.RootUsername, creds.Username, preserveKnown)
	m.SSLConnection = helper.KeepOrUpdateBool(m.SSLConnection, db.SSLConnection, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(db.Status), preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, db.Type, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(db.Updated), preserveKnown)
	m.Version = helper.KeepOrUpdateString(m.Version, db.Version, preserveKnown)

	m.AllowList = helper.KeepOrUpdateSet(
		types.StringType,
		m.AllowList,
		helper.StringSliceToFrameworkValueSlice(db.AllowList),
		preserveKnown,
		&d,
	)
	if d.HasError() {
		return
	}

	membersCasted := helper.MapMap(
		db.Members,
		func(key string, value linodego.DatabaseMemberType) (string, string) {
			return key, string(value)
		},
	)

	m.Members = helper.KeepOrUpdateStringMap(ctx, m.Members, membersCasted, preserveKnown, &d)
	if d.HasError() {
		return
	}

	if db.Fork != nil {
		m.ForkSource = helper.KeepOrUpdateInt64(
			m.ForkSource,
			int64(db.Fork.Source),
			preserveKnown,
		)

		m.ForkRestoreTime = helper.KeepOrUpdateValue(
			m.ForkRestoreTime,
			timetypes.NewRFC3339TimePointerValue(db.Fork.RestoreTime),
			preserveKnown,
		)

	} else {
		m.ForkSource = helper.KeepOrUpdateValue(
			m.ForkSource,
			types.Int64Null(),
			preserveKnown,
		)

		m.ForkRestoreTime = helper.KeepOrUpdateValue(
			m.ForkRestoreTime,
			timetypes.NewRFC3339Null(),
			preserveKnown,
		)
	}

	updatesObject, rd := types.ObjectValueFrom(
		ctx,
		updatesAttributes,
		&ModelUpdates{
			DayOfWeek: types.Int64Value(int64(db.Updates.DayOfWeek)),
			Duration:  types.Int64Value(int64(db.Updates.Duration)),
			Frequency: types.StringValue(string(db.Updates.Frequency)),
			HourOfDay: types.Int64Value(int64(db.Updates.HourOfDay)),
		},
	)
	d.Append(rd...)
	m.Updates = helper.KeepOrUpdateValue(m.Updates, updatesObject, preserveKnown)

	pendingObjects := helper.MapSlice(
		db.Updates.Pending,
		func(pending linodego.DatabaseMaintenanceWindowPending) types.Object {
			result, rd := types.ObjectValueFrom(
				ctx,
				pendingUpdateAttributes,
				&ModelPendingUpdate{
					Deadline:    timetypes.NewRFC3339TimePointerValue(pending.Deadline),
					Description: types.StringValue(pending.Description),
					PlannedFor:  timetypes.NewRFC3339TimePointerValue(pending.PlannedFor),
				},
			)
			d.Append(rd...)

			return result
		},
	)

	pendingSet, rd := types.SetValueFrom(
		ctx,
		types.ObjectType{
			AttrTypes: pendingUpdateAttributes,
		},
		pendingObjects,
	)
	d.Append(rd...)

	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, pendingSet, preserveKnown)

	return nil
}

func (m *Model) CopyFrom(other *Model, preserveKnown bool) {
	m.ForkSource = helper.KeepOrUpdateValue(m.ForkSource, other.ForkSource, preserveKnown)
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.AllowList = helper.KeepOrUpdateValue(m.AllowList, other.AllowList, preserveKnown)
	m.CACert = helper.KeepOrUpdateValue(m.CACert, other.CACert, preserveKnown)
	m.ClusterSize = helper.KeepOrUpdateValue(m.ClusterSize, other.ClusterSize, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Encrypted = helper.KeepOrUpdateValue(m.Encrypted, other.Encrypted, preserveKnown)
	m.Engine = helper.KeepOrUpdateValue(m.Engine, other.Engine, preserveKnown)
	m.EngineID = helper.KeepOrUpdateValue(m.EngineID, other.EngineID, preserveKnown)
	m.ForkRestoreTime = helper.KeepOrUpdateValue(m.ForkRestoreTime, other.ForkRestoreTime, preserveKnown)
	m.HostPrimary = helper.KeepOrUpdateValue(m.HostPrimary, other.HostPrimary, preserveKnown)
	m.HostSecondary = helper.KeepOrUpdateValue(m.HostSecondary, other.HostSecondary, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Members = helper.KeepOrUpdateValue(m.Members, other.Members, preserveKnown)
	m.OldestRestoreTime = helper.KeepOrUpdateValue(m.OldestRestoreTime, other.OldestRestoreTime, preserveKnown)
	m.PendingUpdates = helper.KeepOrUpdateValue(m.PendingUpdates, other.PendingUpdates, preserveKnown)
	m.Platform = helper.KeepOrUpdateValue(m.Platform, other.Platform, preserveKnown)
	m.Port = helper.KeepOrUpdateValue(m.Port, other.Port, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.RootPassword = helper.KeepOrUpdateValue(m.RootPassword, other.RootPassword, preserveKnown)
	m.RootUsername = helper.KeepOrUpdateValue(m.RootUsername, other.RootUsername, preserveKnown)
	m.SSLConnection = helper.KeepOrUpdateValue(m.SSLConnection, other.SSLConnection, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Updates = helper.KeepOrUpdateValue(m.Updates, other.Updates, preserveKnown)
	m.Version = helper.KeepOrUpdateValue(m.Version, other.Version, preserveKnown)
}

// GetFork returns the linodego.DatabaseFork for this model if specified, else nil.
func (m *Model) GetFork(d diag.Diagnostics) *linodego.DatabaseFork {
	var result linodego.DatabaseFork

	isSpecified := false

	if !m.ForkSource.IsUnknown() && !m.ForkSource.IsNull() {
		isSpecified = true

		result.Source = helper.FrameworkSafeInt64ToInt(m.ForkSource.ValueInt64(), &d)
	}

	if !m.ForkRestoreTime.IsUnknown() && !m.ForkRestoreTime.IsNull() {
		isSpecified = true

		restoreTime, rd := m.ForkRestoreTime.ValueRFC3339Time()
		d.Append(rd...)

		result.RestoreTime = &restoreTime
	}

	if d.HasError() || !isSpecified {
		return nil
	}

	return &result
}

// GetAllowList returns the allow list slice for this model if specified, else nil.
func (m *Model) GetAllowList(ctx context.Context, d diag.Diagnostics) []string {
	if m.AllowList.IsUnknown() || m.AllowList.IsNull() {
		return nil
	}

	var result []string

	d.Append(
		m.AllowList.ElementsAs(
			ctx,
			&result,
			false,
		)...,
	)

	return result
}

// GetUpdates returns the ModelUpdates for this model if specified, else nil.
func (m *Model) GetUpdates(ctx context.Context, d diag.Diagnostics) *ModelUpdates {
	if m.Updates.IsUnknown() || m.Updates.IsNull() {
		return nil
	}

	var result ModelUpdates

	d.Append(
		m.Updates.As(
			ctx,
			&result,
			basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true},
		)...,
	)

	return &result
}
