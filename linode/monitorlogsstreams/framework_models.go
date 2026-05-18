package monitorlogsstreams

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// LogsStreamModel is an item in the list data source.
type LogsStreamModel struct {
	ID        types.Int64       `tfsdk:"id"`
	Label     types.String      `tfsdk:"label"`
	Type      types.String      `tfsdk:"type"`
	Status    types.String      `tfsdk:"status"`
	CreatedBy types.String      `tfsdk:"created_by"`
	UpdatedBy types.String      `tfsdk:"updated_by"`
	Created   timetypes.RFC3339 `tfsdk:"created"`
	Updated   timetypes.RFC3339 `tfsdk:"updated"`
	Version   types.Int64       `tfsdk:"version"`
}

// LogsStreamFilterModel is the model for the list data source.
type LogsStreamFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Streams []LogsStreamModel                `tfsdk:"streams"`
}

func (m *LogsStreamFilterModel) ParseLogStreams(streams []linodego.Stream) {
	result := make([]LogsStreamModel, len(streams))
	for i, s := range streams {
		result[i] = LogsStreamModel{
			ID:        types.Int64Value(int64(s.ID)),
			Label:     types.StringValue(s.Label),
			Type:      types.StringValue(string(s.Type)),
			Status:    types.StringValue(string(s.Status)),
			CreatedBy: types.StringValue(s.CreatedBy),
			UpdatedBy: types.StringValue(s.UpdatedBy),
			Created:   timetypes.NewRFC3339TimePointerValue(s.Created),
			Updated:   timetypes.NewRFC3339TimePointerValue(s.Updated),
			Version:   types.Int64Value(int64(s.Version)),
		}
	}
	m.Streams = result
}

var streamDetailsAttrTypes = map[string]attr.Type{
	"cluster_ids":                      types.ListType{ElemType: types.Int64Type},
	"is_auto_add_all_clusters_enabled": types.BoolType,
}

// StreamDetailsModel holds the details block for an lke_audit_logs stream.
type StreamDetailsModel struct {
	ClusterIDs                  types.List `tfsdk:"cluster_ids"`
	IsAutoAddAllClustersEnabled types.Bool `tfsdk:"is_auto_add_all_clusters_enabled"`
}

// ResourceModel represents the linode_monitor_logs_stream resource state.
type ResourceModel struct {
	ID           types.String      `tfsdk:"id"`
	Label        types.String      `tfsdk:"label"`
	Type         types.String      `tfsdk:"type"`
	Status       types.String      `tfsdk:"status"`
	Version      types.Int64       `tfsdk:"version"`
	Destinations types.List        `tfsdk:"destinations"`
	Details      types.Object      `tfsdk:"details"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	CreatedBy    types.String      `tfsdk:"created_by"`
	UpdatedBy    types.String      `tfsdk:"updated_by"`
}

// DataSourceModel represents the linode_monitor_logs_stream data source state.
type DataSourceModel struct {
	ID           types.String      `tfsdk:"id"`
	Label        types.String      `tfsdk:"label"`
	Type         types.String      `tfsdk:"type"`
	Status       types.String      `tfsdk:"status"`
	Version      types.Int64       `tfsdk:"version"`
	Destinations types.List        `tfsdk:"destinations"`
	Details      types.Object      `tfsdk:"details"`
	Created      timetypes.RFC3339 `tfsdk:"created"`
	Updated      timetypes.RFC3339 `tfsdk:"updated"`
	CreatedBy    types.String      `tfsdk:"created_by"`
	UpdatedBy    types.String      `tfsdk:"updated_by"`
}

// FlattenStream populates the ResourceModel from a linodego.Stream.
func (m *ResourceModel) FlattenStream(
	ctx context.Context,
	stream *linodego.Stream,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(stream.ID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, stream.Label, preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(stream.Type), preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(stream.Status), preserveKnown)
	m.Version = helper.KeepOrUpdateInt64(m.Version, int64(stream.Version), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(
		m.Created,
		timetypes.NewRFC3339TimePointerValue(stream.Created),
		preserveKnown,
	)
	m.Updated = helper.KeepOrUpdateValue(
		m.Updated,
		timetypes.NewRFC3339TimePointerValue(stream.Updated),
		preserveKnown,
	)
	m.CreatedBy = helper.KeepOrUpdateString(m.CreatedBy, stream.CreatedBy, preserveKnown)
	m.UpdatedBy = helper.KeepOrUpdateString(m.UpdatedBy, stream.UpdatedBy, preserveKnown)

	destIDs := helper.MapSlice(stream.Destinations, func(d linodego.StreamDestination) int64 {
		return int64(d.ID)
	})
	destList, d := types.ListValueFrom(ctx, types.Int64Type, destIDs)
	diags.Append(d...)
	if !diags.HasError() {
		m.Destinations = helper.KeepOrUpdateValue(m.Destinations, destList, preserveKnown)
	}

	m.Details = flattenStreamDetails(ctx, stream.Details, m.Details, preserveKnown, diags)
}

// FlattenStream populates the DataSourceModel from a linodego.Stream.
func (m *DataSourceModel) FlattenStream(
	ctx context.Context,
	stream *linodego.Stream,
	diags *diag.Diagnostics,
) {
	rm := ResourceModel{
		Details:      m.Details,
		Destinations: m.Destinations,
	}
	rm.FlattenStream(ctx, stream, false, diags)
	m.ID = rm.ID
	m.Label = rm.Label
	m.Type = rm.Type
	m.Status = rm.Status
	m.Version = rm.Version
	m.Destinations = rm.Destinations
	m.Details = rm.Details
	m.Created = rm.Created
	m.Updated = rm.Updated
	m.CreatedBy = rm.CreatedBy
	m.UpdatedBy = rm.UpdatedBy
}

// GetCreateOptions builds linodego.StreamCreateOptions from the ResourceModel.
func (m *ResourceModel) GetCreateOptions(
	ctx context.Context,
	diags *diag.Diagnostics,
) linodego.StreamCreateOptions {
	opts := linodego.StreamCreateOptions{
		Label: m.Label.ValueString(),
		Type:  linodego.StreamType(m.Type.ValueString()),
	}

	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		status := linodego.StreamStatus(m.Status.ValueString())
		opts.Status = &status
	}

	var destIDs []int64
	diags.Append(m.Destinations.ElementsAs(ctx, &destIDs, false)...)
	if !diags.HasError() {
		opts.Destinations = helper.MapSlice(destIDs, func(id int64) int { return int(id) })
	}

	if !m.Details.IsNull() && !m.Details.IsUnknown() {
		var detailsModel StreamDetailsModel
		diags.Append(m.Details.As(ctx, &detailsModel, basetypes.ObjectAsOptions{})...)
		if !diags.HasError() {
			opts.Details = buildStreamDetailsOptions(ctx, detailsModel, diags)
		}
	}

	return opts
}

// GetUpdateOptions builds linodego.StreamUpdateOptions from the ResourceModel.
func (m *ResourceModel) GetUpdateOptions(
	ctx context.Context,
	diags *diag.Diagnostics,
) linodego.StreamUpdateOptions {
	opts := linodego.StreamUpdateOptions{
		Label: linodego.Pointer(m.Label.ValueString()),
	}

	streamType := linodego.StreamType(m.Type.ValueString())
	opts.Type = &streamType

	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		status := linodego.StreamStatus(m.Status.ValueString())
		opts.Status = &status
	}

	if !m.Destinations.IsNull() && !m.Destinations.IsUnknown() {
		var destIDs []int64
		diags.Append(m.Destinations.ElementsAs(ctx, &destIDs, false)...)
		if !diags.HasError() {
			opts.Destinations = helper.MapSlice(destIDs, func(id int64) int { return int(id) })
		}
	}

	if !m.Details.IsNull() && !m.Details.IsUnknown() {
		var detailsModel StreamDetailsModel
		diags.Append(m.Details.As(ctx, &detailsModel, basetypes.ObjectAsOptions{})...)
		if !diags.HasError() {
			opts.Details = buildStreamDetailsOptions(ctx, detailsModel, diags)
		}
	}

	return opts
}

func flattenStreamDetails(
	ctx context.Context,
	details *linodego.StreamDetails,
	original types.Object,
	preserveKnown bool,
	diags *diag.Diagnostics,
) types.Object {
	if details == nil {
		if preserveKnown && !original.IsNull() {
			return original
		}
		return types.ObjectNull(streamDetailsAttrTypes)
	}

	return *helper.KeepOrUpdateSingleNestedAttributesWithTypes(
		ctx, original, streamDetailsAttrTypes, preserveKnown, diags,
		func(model *StreamDetailsModel, isNull *bool, preserveKnown bool, diags *diag.Diagnostics) {
			clusterIDs := helper.MapSlice(details.ClusterIDs, func(id int) int64 { return int64(id) })
			clusterList, d := types.ListValueFrom(ctx, types.Int64Type, clusterIDs)
			diags.Append(d...)
			if !diags.HasError() {
				model.ClusterIDs = helper.KeepOrUpdateValue(model.ClusterIDs, clusterList, preserveKnown)
			}
			model.IsAutoAddAllClustersEnabled = helper.KeepOrUpdateBool(
				model.IsAutoAddAllClustersEnabled,
				details.IsAutoAddAllClustersEnabled,
				preserveKnown,
			)
		},
	)
}

func buildStreamDetailsOptions(
	ctx context.Context,
	model StreamDetailsModel,
	diags *diag.Diagnostics,
) *linodego.StreamDetails {
	details := &linodego.StreamDetails{
		IsAutoAddAllClustersEnabled: model.IsAutoAddAllClustersEnabled.ValueBool(),
	}

	if !model.ClusterIDs.IsNull() && !model.ClusterIDs.IsUnknown() {
		var clusterIDs []int64
		diags.Append(model.ClusterIDs.ElementsAs(ctx, &clusterIDs, false)...)
		if !diags.HasError() {
			details.ClusterIDs = helper.MapSlice(clusterIDs, func(id int64) int { return int(id) })
		}
	}

	return details
}
