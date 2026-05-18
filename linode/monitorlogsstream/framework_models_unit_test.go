//go:build unit

package monitorlogsstream_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTime(s string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		panic(err)
	}
	return &t
}

// streamDetailsAttrTypes derives the details nested object attribute types from the exported schema.
func streamDetailsAttrTypes() map[string]attr.Type {
	detailsAttr := monitorlogsstream.ResourceSchema.Attributes["details"].(schema.SingleNestedAttribute)
	return detailsAttr.GetType().(types.ObjectType).AttrTypes
}

func TestFlattenStream_AuditLogs(t *testing.T) {
	ctx := context.Background()

	destID := 12345
	stream := &linodego.Stream{
		ID:      456,
		Label:   "AuditLog-config",
		Type:    linodego.StreamTypeAuditLogs,
		Status:  linodego.StreamStatusActive,
		Version: 1,
		Destinations: []linodego.StreamDestination{
			{ID: destID, Label: "OBJ_logs_destination", Type: linodego.StreamDestinationTypeAkamaiObjectStorage},
		},
		Details:   nil,
		Created:   makeTime("2025-03-20T01:41:09"),
		Updated:   makeTime("2025-03-20T01:41:09"),
		CreatedBy: "John Q. Linode",
		UpdatedBy: "Jane Q. Linode",
	}

	var m monitorlogsstream.ResourceModel
	var diags diag.Diagnostics

	m.FlattenStream(ctx, stream, false, &diags)

	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "456", m.ID.ValueString())
	assert.Equal(t, "AuditLog-config", m.Label.ValueString())
	assert.Equal(t, "audit_logs", m.Type.ValueString())
	assert.Equal(t, "active", m.Status.ValueString())
	assert.Equal(t, int64(1), m.Version.ValueInt64())
	assert.Equal(t, "John Q. Linode", m.CreatedBy.ValueString())
	assert.Equal(t, "Jane Q. Linode", m.UpdatedBy.ValueString())
	assert.False(t, m.Created.IsNull())
	assert.False(t, m.Updated.IsNull())

	// details should be null for audit_logs
	assert.True(t, m.Details.IsNull(), "details should be null for audit_logs stream")
	_ = streamDetailsAttrTypes() // ensure schema is accessible

	// destinations should contain the destination ID
	var destIDs []int64
	require.False(t, m.Destinations.IsNull())
	diags = m.Destinations.ElementsAs(ctx, &destIDs, false)
	require.False(t, diags.HasError())
	assert.Equal(t, []int64{int64(destID)}, destIDs)
}

func TestFlattenStream_LKEAuditLogs(t *testing.T) {
	ctx := context.Background()

	stream := &linodego.Stream{
		ID:      789,
		Label:   "LKEAuditLog-config",
		Type:    linodego.StreamTypeLKEAuditLogs,
		Status:  linodego.StreamStatusActive,
		Version: 1,
		Destinations: []linodego.StreamDestination{
			{ID: 12345, Label: "OBJ_logs_destination", Type: linodego.StreamDestinationTypeAkamaiObjectStorage},
		},
		Details: &linodego.StreamDetails{
			ClusterIDs:                  []int{1234, 5678},
			IsAutoAddAllClustersEnabled: false,
		},
		Created:   makeTime("2025-03-20T01:41:09"),
		Updated:   makeTime("2025-03-20T01:41:09"),
		CreatedBy: "John Q. Linode",
		UpdatedBy: "Jane Q. Linode",
	}

	var m monitorlogsstream.ResourceModel
	var diags diag.Diagnostics

	m.FlattenStream(ctx, stream, false, &diags)

	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

	assert.Equal(t, "789", m.ID.ValueString())
	assert.Equal(t, "lke_audit_logs", m.Type.ValueString())
	assert.False(t, m.Details.IsNull(), "details should be set for lke_audit_logs stream")

	var detailsModel monitorlogsstream.StreamDetailsModel
	diags = m.Details.As(ctx, &detailsModel, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError())

	assert.False(t, detailsModel.IsAutoAddAllClustersEnabled.ValueBool())

	var clusterIDs []int64
	diags = detailsModel.ClusterIDs.ElementsAs(ctx, &clusterIDs, false)
	require.False(t, diags.HasError())
	assert.Equal(t, []int64{1234, 5678}, clusterIDs)
}

func TestFlattenStream_PreserveKnown(t *testing.T) {
	ctx := context.Background()

	stream := &linodego.Stream{
		ID:    456,
		Label: "AuditLog-updated",
		Type:  linodego.StreamTypeAuditLogs,
		Destinations: []linodego.StreamDestination{
			{ID: 99999},
		},
		Created:   makeTime("2025-03-20T01:41:09"),
		Updated:   makeTime("2025-03-20T01:41:09"),
		CreatedBy: "user",
		UpdatedBy: "user",
	}

	// Pre-populate with known values that should be preserved
	m := monitorlogsstream.ResourceModel{
		ID:      types.StringValue("456"),
		Label:   types.StringValue("AuditLog-original"),
		Details: types.ObjectNull(streamDetailsAttrTypes()),
	}
	var diags diag.Diagnostics

	m.FlattenStream(ctx, stream, true, &diags)
	require.False(t, diags.HasError())

	// Known values should be preserved when preserveKnown = true
	assert.Equal(t, "AuditLog-original", m.Label.ValueString(), "known label should be preserved")
	assert.Equal(t, "456", m.ID.ValueString(), "known ID should be preserved")
}

func TestFlattenStream_NotPreserveKnown(t *testing.T) {
	ctx := context.Background()

	stream := &linodego.Stream{
		ID:    456,
		Label: "AuditLog-updated",
		Type:  linodego.StreamTypeAuditLogs,
		Destinations: []linodego.StreamDestination{
			{ID: 99999},
		},
		Created:   makeTime("2025-03-20T01:41:09"),
		Updated:   makeTime("2025-03-20T01:41:09"),
		CreatedBy: "user",
		UpdatedBy: "user",
	}

	m := monitorlogsstream.ResourceModel{
		ID:      types.StringValue("456"),
		Label:   types.StringValue("AuditLog-original"),
		Details: types.ObjectNull(streamDetailsAttrTypes()),
	}
	var diags diag.Diagnostics

	m.FlattenStream(ctx, stream, false, &diags)
	require.False(t, diags.HasError())

	// Values should be updated when preserveKnown = false
	assert.Equal(t, "AuditLog-updated", m.Label.ValueString(), "label should be updated")
}

func TestGetCreateOptions_AuditLogs(t *testing.T) {
	ctx := context.Background()

	destList, _ := types.ListValueFrom(ctx, types.Int64Type, []int64{12345})

	m := monitorlogsstream.ResourceModel{
		Label:        types.StringValue("AuditLog-config"),
		Type:         types.StringValue("audit_logs"),
		Destinations: destList,
		Details:      types.ObjectNull(streamDetailsAttrTypes()),
	}

	var diags diag.Diagnostics
	opts := m.GetCreateOptions(ctx, &diags)
	require.False(t, diags.HasError())

	assert.Equal(t, "AuditLog-config", opts.Label)
	assert.Equal(t, linodego.StreamTypeAuditLogs, opts.Type)
	assert.Equal(t, []int{12345}, opts.Destinations)
	assert.Nil(t, opts.Details)
	assert.Nil(t, opts.Status)
}

func TestGetCreateOptions_LKEAuditLogs(t *testing.T) {
	ctx := context.Background()

	destList, _ := types.ListValueFrom(ctx, types.Int64Type, []int64{12345})
	clusterList, _ := types.ListValueFrom(ctx, types.Int64Type, []int64{1111, 2222})

	detailsObj, _ := types.ObjectValueFrom(ctx, streamDetailsAttrTypes(), monitorlogsstream.StreamDetailsModel{
		ClusterIDs:                  clusterList,
		IsAutoAddAllClustersEnabled: types.BoolValue(false),
	})

	m := monitorlogsstream.ResourceModel{
		Label:        types.StringValue("LKEAuditLog-config"),
		Type:         types.StringValue("lke_audit_logs"),
		Destinations: destList,
		Details:      detailsObj,
	}

	var diags diag.Diagnostics
	opts := m.GetCreateOptions(ctx, &diags)
	require.False(t, diags.HasError())

	assert.Equal(t, "LKEAuditLog-config", opts.Label)
	assert.Equal(t, linodego.StreamTypeLKEAuditLogs, opts.Type)
	require.NotNil(t, opts.Details)
	assert.Equal(t, []int{1111, 2222}, opts.Details.ClusterIDs)
	assert.False(t, opts.Details.IsAutoAddAllClustersEnabled)
}

func TestGetCreateOptions_WithStatus(t *testing.T) {
	ctx := context.Background()

	destList, _ := types.ListValueFrom(ctx, types.Int64Type, []int64{12345})

	m := monitorlogsstream.ResourceModel{
		Label:        types.StringValue("AuditLog-config"),
		Type:         types.StringValue("audit_logs"),
		Status:       types.StringValue("inactive"),
		Destinations: destList,
		Details:      types.ObjectNull(streamDetailsAttrTypes()),
	}

	var diags diag.Diagnostics
	opts := m.GetCreateOptions(ctx, &diags)
	require.False(t, diags.HasError())

	require.NotNil(t, opts.Status)
	assert.Equal(t, linodego.StreamStatusInactive, *opts.Status)
}

func TestGetUpdateOptions(t *testing.T) {
	ctx := context.Background()

	destList, _ := types.ListValueFrom(ctx, types.Int64Type, []int64{12345})

	m := monitorlogsstream.ResourceModel{
		Label:        types.StringValue("AuditLog-renamed"),
		Type:         types.StringValue("audit_logs"),
		Status:       types.StringValue("active"),
		Destinations: destList,
		Details:      types.ObjectNull(streamDetailsAttrTypes()),
	}

	var diags diag.Diagnostics
	opts := m.GetUpdateOptions(ctx, &diags)
	require.False(t, diags.HasError())

	require.NotNil(t, opts.Label)
	assert.Equal(t, "AuditLog-renamed", *opts.Label)
	assert.Equal(t, []int{12345}, opts.Destinations)
	require.NotNil(t, opts.Status)
	assert.Equal(t, linodego.StreamStatusActive, *opts.Status)
}
