package monitorlogsstreamhistory

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

var streamDetailsAttrTypes = map[string]attr.Type{
	"cluster_ids":                      types.ListType{ElemType: types.Int64Type},
	"is_auto_add_all_clusters_enabled": types.BoolType,
}

// StreamHistoryEntryModel represents a single historical version of a stream.
type StreamHistoryEntryModel struct {
	ID           types.Int64       `tfsdk:"id"`
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

// DataSourceModel represents the linode_monitor_logs_stream_history data source state.
type DataSourceModel struct {
	ID       types.String              `tfsdk:"id"`
	StreamID types.Int64               `tfsdk:"stream_id"`
	Streams  []StreamHistoryEntryModel `tfsdk:"streams"`
}

func flattenStreamHistoryEntry(
	ctx context.Context,
	stream linodego.Stream,
	diags *diag.Diagnostics,
) StreamHistoryEntryModel {
	entry := StreamHistoryEntryModel{
		ID:        types.Int64Value(int64(stream.ID)),
		Label:     types.StringValue(stream.Label),
		Type:      types.StringValue(string(stream.Type)),
		Status:    types.StringValue(string(stream.Status)),
		Version:   types.Int64Value(int64(stream.Version)),
		Created:   timetypes.NewRFC3339TimePointerValue(stream.Created),
		Updated:   timetypes.NewRFC3339TimePointerValue(stream.Updated),
		CreatedBy: types.StringValue(stream.CreatedBy),
		UpdatedBy: types.StringValue(stream.UpdatedBy),
	}

	destIDs := helper.MapSlice(stream.Destinations, func(d linodego.StreamDestination) int64 {
		return int64(d.ID)
	})
	destList, d := types.ListValueFrom(ctx, types.Int64Type, destIDs)
	diags.Append(d...)
	if !diags.HasError() {
		entry.Destinations = destList
	}

	if stream.Details != nil {
		clusterIDs := helper.MapSlice(stream.Details.ClusterIDs, func(id int) int64 { return int64(id) })
		clusterList, d := types.ListValueFrom(ctx, types.Int64Type, clusterIDs)
		diags.Append(d...)

		if diags.HasError() {
			entry.Details = types.ObjectNull(streamDetailsAttrTypes)
		} else {
			detailsObj, d := types.ObjectValue(streamDetailsAttrTypes, map[string]attr.Value{
				"cluster_ids":                      clusterList,
				"is_auto_add_all_clusters_enabled": types.BoolValue(stream.Details.IsAutoAddAllClustersEnabled),
			})
			diags.Append(d...)
			entry.Details = detailsObj
		}
	} else {
		entry.Details = types.ObjectNull(streamDetailsAttrTypes)
	}

	return entry
}

func (m *DataSourceModel) ParseStreamHistory(
	ctx context.Context,
	streams []linodego.Stream,
	diags *diag.Diagnostics,
) {
	m.ID = types.StringValue(strconv.FormatInt(m.StreamID.ValueInt64(), 10))
	result := make([]StreamHistoryEntryModel, len(streams))
	for i, s := range streams {
		result[i] = flattenStreamHistoryEntry(ctx, s, diags)
		if diags.HasError() {
			return
		}
	}
	m.Streams = result
}
