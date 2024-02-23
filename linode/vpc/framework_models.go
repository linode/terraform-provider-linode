package vpc

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type VPCModel struct {
	ID          types.String      `tfsdk:"id"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	Region      types.String      `tfsdk:"region"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
}

func (m *VPCModel) FlattenVPC(ctx context.Context, vpc *linodego.VPC, preserveKnown bool) {
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
}

func (m *VPCModel) CopyFrom(ctx context.Context, other VPCModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.Description = helper.KeepOrUpdateValue(m.Description, other.Description, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
}
