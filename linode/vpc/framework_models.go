package vpc

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type VPCModel struct {
	ID          types.String      `tfsdk:"id"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	Region      types.String      `tfsdk:"region"`
	IPv6        types.Set         `tfsdk:"ipv6"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
}

type VPCIPv6Model struct {
	Range           types.String `tfsdk:"range"`
	AllocationClass types.String `tfsdk:"allocation_class"`
}

var VPCIPv6ModelObjectType = helper.Must(
	helper.FrameworkModelToObjectType[VPCIPv6Model](context.Background()),
)

func (m *VPCModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) diag.Diagnostics {
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

	ipv6Set, diags := types.SetValueFrom(ctx, VPCIPv6ModelObjectType, vpc.IPv6)
	if diags.HasError() {
		return diags
	}

	m.IPv6 = helper.KeepOrUpdateValue(
		m.IPv6,
		ipv6Set,
		preserveKnown,
	)

	return nil
}

func (m *VPCModel) CopyFrom(ctx context.Context, other VPCModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.Description = helper.KeepOrUpdateValue(m.Description, other.Description, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.IPv6 = helper.KeepOrUpdateValue(m.IPv6, other.IPv6, preserveKnown)
}
