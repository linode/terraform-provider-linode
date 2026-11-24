package nfsspace

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// BaseModel contains the fields shared between the resource and data source models.
type BaseModel struct {
	ID                types.String      `tfsdk:"id"`
	Label             types.String      `tfsdk:"label"`
	Description       types.String      `tfsdk:"description"`
	Status            types.String      `tfsdk:"status"`
	MaxCapacityBytes  types.Int64       `tfsdk:"max_capacity_bytes"`
	UsedCapacityBytes types.Int64       `tfsdk:"used_capacity_bytes"`
	Created           timetypes.RFC3339 `tfsdk:"created"`
	Updated           timetypes.RFC3339 `tfsdk:"updated"`
	Tags              types.Set         `tfsdk:"tags"`
}

func (m *BaseModel) Flatten(
	ctx context.Context,
	space *linodego.NFSSpace,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(space.ID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, space.Label, preserveKnown)
	m.Description = helper.KeepOrUpdateStringPointer(m.Description, space.Description, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(space.Status), preserveKnown)
	m.MaxCapacityBytes = helper.KeepOrUpdateInt64(m.MaxCapacityBytes, space.MaxCapacityBytes, preserveKnown)
	m.UsedCapacityBytes = helper.KeepOrUpdateInt64(m.UsedCapacityBytes, space.UsedCapacityBytes, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(
		m.Created,
		timetypes.NewRFC3339TimePointerValue(space.Created),
		preserveKnown,
	)
	m.Updated = helper.KeepOrUpdateValue(
		m.Updated,
		timetypes.NewRFC3339TimePointerValue(space.Updated),
		preserveKnown,
	)

	tags, d := types.SetValueFrom(ctx, types.StringType, space.Tags)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Tags = helper.KeepOrUpdateValue(m.Tags, tags, preserveKnown)

	return diags
}

func (m *BaseModel) CopyFrom(other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Description = helper.KeepOrUpdateValue(m.Description, other.Description, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.MaxCapacityBytes = helper.KeepOrUpdateValue(m.MaxCapacityBytes, other.MaxCapacityBytes, preserveKnown)
	m.UsedCapacityBytes = helper.KeepOrUpdateValue(m.UsedCapacityBytes, other.UsedCapacityBytes, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Tags = helper.KeepOrUpdateValue(m.Tags, other.Tags, preserveKnown)
}

// ResourceModel represents the Terraform state for the linode_nfs_space resource.
type ResourceModel struct {
	BaseModel
}

// DataSourceModel represents the Terraform state for the linode_nfs_space data source.
type DataSourceModel struct {
	BaseModel
}
