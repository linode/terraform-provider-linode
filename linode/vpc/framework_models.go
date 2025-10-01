package vpc

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/customtypes"
)

/*
Shared Implementation
*/

type BaseModel struct {
	ID          types.String      `tfsdk:"id"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	Region      types.String      `tfsdk:"region"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
}

func (m *BaseModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(vpc.ID), preserveKnown)

	m.Description = helper.KeepOrUpdateString(m.Description, vpc.Description, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(
		m.Created,
		timetypes.NewRFC3339TimePointerValue(vpc.Created),
		preserveKnown,
	)
	m.Updated = helper.KeepOrUpdateValue(
		m.Updated,
		timetypes.NewRFC3339TimePointerValue(vpc.Updated),
		preserveKnown,
	)
	m.Label = helper.KeepOrUpdateString(m.Label, vpc.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, vpc.Region, preserveKnown)

	return nil
}

func (m *BaseModel) CopyFrom(ctx context.Context, other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.Description = helper.KeepOrUpdateValue(m.Description, other.Description, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
}

/*
Resource-Specific Implementation
*/

type ResourceModel struct {
	BaseModel
	IPv6 types.List `tfsdk:"ipv6"`
}

type ResourceModelIPv6 struct {
	Range           customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
	AllocatedRange  types.String                          `tfsdk:"allocated_range"`
	AllocationClass types.String                          `tfsdk:"allocation_class"`
}

func (m *ResourceModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
	m.BaseModel.FlattenVPC(ctx, vpc, preserveKnown)

	ipv6Models := helper.MapSlice(vpc.IPv6,
		func(r linodego.VPCIPv6Range) ResourceModelIPv6 {
			return ResourceModelIPv6{
				Range:          customtypes.LinodeAutoAllocRangeValue{StringValue: types.StringValue(r.Range)},
				AllocatedRange: types.StringValue(r.Range),
			}
		},
	)

	ipv6List, diags := types.ListValueFrom(ctx, ResourceSchemaIPv6NestedObject.Type(), ipv6Models)
	if diags.HasError() {
		return diags
	}

	m.IPv6 = helper.KeepOrUpdateValue(
		m.IPv6,
		ipv6List,
		// NOTE: preserveKnown is false here to ensure the allocated_range attribute is populated
		false,
	)

	return nil
}

func (m *ResourceModel) CopyFrom(ctx context.Context, other ResourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(ctx, other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
}

/*
Data Source-Specific Implementation
*/

type DataSourceModel struct {
	BaseModel
	IPv6 types.List `tfsdk:"ipv6"`
}

type DataSourceModelIPv6 struct {
	Range customtypes.LinodeAutoAllocRangeValue `tfsdk:"range"`
}

func (m *DataSourceModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
	m.BaseModel.FlattenVPC(ctx, vpc, preserveKnown)

	ipv6Models := helper.MapSlice(
		vpc.IPv6,
		func(r linodego.VPCIPv6Range) DataSourceModelIPv6 {
			return DataSourceModelIPv6{
				Range: customtypes.LinodeAutoAllocRangeValue{StringValue: types.StringValue(r.Range)},
			}
		},
	)

	ipv6List, diags := types.ListValueFrom(ctx, DataSourceSchemaIPv6NestedObject.Type(), ipv6Models)
	if diags.HasError() {
		return diags
	}

	m.IPv6 = helper.KeepOrUpdateValue(
		m.IPv6,
		ipv6List,
		preserveKnown,
	)

	return nil
}

func (m *DataSourceModel) CopyFrom(ctx context.Context, other DataSourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(ctx, other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
}
