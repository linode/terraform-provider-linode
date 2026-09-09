package vpc

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/customtypes"
	"github.com/linode/terraform-provider-linode/v4/linode/vpcsubnet"
)

/*
Shared Implementation
*/

type BaseModel struct {
	ID          types.String      `tfsdk:"id"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	Region      types.String      `tfsdk:"region"`
	VPCType     types.String      `tfsdk:"vpc_type"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
	IPv4        types.List        `tfsdk:"ipv4"`
	Subnets     types.List        `tfsdk:"subnets"`
}

type SubnetModel struct {
	ID            types.Int64       `tfsdk:"id"`
	Label         types.String      `tfsdk:"label"`
	IPv4          types.String      `tfsdk:"ipv4"`
	IPv6          types.List        `tfsdk:"ipv6"`
	Linodes       types.List        `tfsdk:"linodes"`
	Databases     types.List        `tfsdk:"databases"`
	Nodebalancers types.List        `tfsdk:"nodebalancers"`
	Created       timetypes.RFC3339 `tfsdk:"created"`
	Updated       timetypes.RFC3339 `tfsdk:"updated"`
}

type SubnetModelIPv6 struct {
	Range types.String `tfsdk:"range"`
}

func FlattenSubnets(
	ctx context.Context,
	subnets []linodego.VPCSubnet,
	subnetObjectType attr.Type,
	ipv6ObjectType attr.Type,
	preserveKnown bool,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	subnetModels := make([]SubnetModel, len(subnets))
	for i, subnet := range subnets {
		var model SubnetModel

		model.ID = helper.KeepOrUpdateInt64(model.ID, int64(subnet.ID), preserveKnown)
		model.Label = helper.KeepOrUpdateString(model.Label, subnet.Label, preserveKnown)
		model.IPv4 = helper.KeepOrUpdateString(model.IPv4, subnet.IPv4, preserveKnown)
		model.Created = helper.KeepOrUpdateValue(
			model.Created,
			timetypes.NewRFC3339TimePointerValue(subnet.Created),
			preserveKnown,
		)
		model.Updated = helper.KeepOrUpdateValue(
			model.Updated,
			timetypes.NewRFC3339TimePointerValue(subnet.Updated),
			preserveKnown,
		)

		ipv6Models := helper.MapSlice(
			subnet.IPv6,
			func(r linodego.VPCIPv6Range) SubnetModelIPv6 {
				return SubnetModelIPv6{
					Range: types.StringValue(r.Range),
				}
			},
		)

		ipv6List, d := types.ListValueFrom(ctx, ipv6ObjectType, ipv6Models)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(subnetObjectType), diags
		}
		model.IPv6 = helper.KeepOrUpdateValue(model.IPv6, ipv6List, preserveKnown)

		linodesList, d := vpcsubnet.FlattenSubnetLinodes(ctx, subnet.Linodes)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(subnetObjectType), diags
		}
		model.Linodes = helper.KeepOrUpdateValue(model.Linodes, *linodesList, preserveKnown)

		databasesList, d := vpcsubnet.FlattenSubnetDatabases(ctx, subnet.Databases)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(subnetObjectType), diags
		}
		model.Databases = helper.KeepOrUpdateValue(model.Databases, *databasesList, preserveKnown)

		nodebalancersList, d := vpcsubnet.FlattenSubnetNodebalancers(ctx, subnet.Nodebalancers)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(subnetObjectType), diags
		}
		model.Nodebalancers = helper.KeepOrUpdateValue(model.Nodebalancers, *nodebalancersList, preserveKnown)

		subnetModels[i] = model
	}

	subnetList, d := types.ListValueFrom(ctx, subnetObjectType, subnetModels)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(subnetObjectType), diags
	}

	return subnetList, diags
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
	m.VPCType = helper.KeepOrUpdateString(m.VPCType, string(vpc.VPCType), preserveKnown)

	return nil
}

func (m *BaseModel) CopyFrom(ctx context.Context, other BaseModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.Description = helper.KeepOrUpdateValue(m.Description, other.Description, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.VPCType = helper.KeepOrUpdateValue(m.VPCType, other.VPCType, preserveKnown)
	m.Subnets = helper.KeepOrUpdateValue(m.Subnets, other.Subnets, preserveKnown)
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

type ResourceModelIPv4 struct {
	Range types.String `tfsdk:"range"`
}

func (m *ResourceModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
	m.BaseModel.FlattenVPC(ctx, vpc, preserveKnown)

	// Flatten IPv6
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

	// Flatten IPv4
	ipv4Models := helper.MapSlice(vpc.IPv4,
		func(r linodego.VPCIPv4Range) ResourceModelIPv4 {
			return ResourceModelIPv4{
				Range: types.StringValue(r.Range),
			}
		},
	)

	ipv4List, ipv4Diags := types.ListValueFrom(ctx, ResourceSchemaIPv4NestedObject.Type(), ipv4Models)
	if ipv4Diags.HasError() {
		return ipv4Diags
	}

	m.IPv4 = helper.KeepOrUpdateValue(
		m.IPv4,
		ipv4List,
		false,
	)

	// Flatten Subnets
	// NOTE: preserveKnown is false here because subnets is computed-only and must
	// always be populated from the API response.
	subnetList, subnetDiags := FlattenSubnets(
		ctx,
		vpc.Subnets,
		ResourceSchemaSubnetNestedObject.Type(),
		ResourceSchemaSubnetIPv6NestedObject.Type(),
		false,
	)
	if subnetDiags.HasError() {
		return subnetDiags
	}

	m.Subnets = helper.KeepOrUpdateValue(m.Subnets, subnetList, false)

	return nil
}

func (m *ResourceModel) CopyFrom(ctx context.Context, other ResourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(ctx, other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
	m.IPv4 = helper.KeepOrUpdateValue(m.IPv4, other.IPv4, preserveKnown)
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

type DataSourceModelIPv4 struct {
	Range types.String `tfsdk:"range"`
}

func (m *DataSourceModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
	m.BaseModel.FlattenVPC(ctx, vpc, preserveKnown)

	// Flatten IPv6
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

	// Flatten IPv4
	ipv4Models := helper.MapSlice(
		vpc.IPv4,
		func(r linodego.VPCIPv4Range) DataSourceModelIPv4 {
			return DataSourceModelIPv4{
				Range: types.StringValue(r.Range),
			}
		},
	)

	ipv4List, ipv4Diags := types.ListValueFrom(ctx, DataSourceSchemaIPv4NestedObject.Type(), ipv4Models)
	if ipv4Diags.HasError() {
		return ipv4Diags
	}

	m.IPv4 = helper.KeepOrUpdateValue(
		m.IPv4,
		ipv4List,
		preserveKnown,
	)

	// Flatten Subnets
	subnetList, subnetDiags := FlattenSubnets(
		ctx,
		vpc.Subnets,
		DataSourceSchemaSubnetNestedObject.Type(),
		DataSourceSchemaSubnetIPv6NestedObject.Type(),
		preserveKnown,
	)
	if subnetDiags.HasError() {
		return subnetDiags
	}

	m.Subnets = helper.KeepOrUpdateValue(m.Subnets, subnetList, preserveKnown)

	return nil
}

func (m *DataSourceModel) CopyFrom(ctx context.Context, other DataSourceModel, preserveKnown bool) {
	m.BaseModel.CopyFrom(ctx, other.BaseModel, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
	m.IPv4 = helper.KeepOrUpdateValue(m.IPv4, other.IPv4, preserveKnown)
}
