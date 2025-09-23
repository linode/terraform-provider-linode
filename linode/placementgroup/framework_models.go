package placementgroup

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type PlacementGroupDataSourceModel struct {
	ID                   types.Int64                   `tfsdk:"id"`
	Label                types.String                  `tfsdk:"label"`
	Region               types.String                  `tfsdk:"region"`
	PlacementGroupType   types.String                  `tfsdk:"placement_group_type"`
	IsCompliant          types.Bool                    `tfsdk:"is_compliant"`
	PlacementGroupPolicy types.String                  `tfsdk:"placement_group_policy"`
	Members              []PlacementGroupMemberModel   `tfsdk:"members"`
	Migrations           *PlacementGroupMigrationModel `tfsdk:"migrations"`
}

type PlacementGroupResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Label                types.String `tfsdk:"label"`
	Region               types.String `tfsdk:"region"`
	PlacementGroupType   types.String `tfsdk:"placement_group_type"`
	PlacementGroupPolicy types.String `tfsdk:"placement_group_policy"`
	IsCompliant          types.Bool   `tfsdk:"is_compliant"`

	Members types.Set `tfsdk:"members"`
}

type PlacementGroupMigrationModel struct {
	Inbound  []PlacementGroupMigrationInstanceModel `tfsdk:"inbound"`
	Outbound []PlacementGroupMigrationInstanceModel `tfsdk:"outbound"`
}

type PlacementGroupMigrationInstanceModel struct {
	LinodeID types.Int64 `tfsdk:"linode_id"`
}

type PlacementGroupMemberModel struct {
	LinodeID    types.Int64 `tfsdk:"linode_id"`
	IsCompliant types.Bool  `tfsdk:"is_compliant"`
}

func (data *PlacementGroupDataSourceModel) ParsePlacementGroup(
	pg *linodego.PlacementGroup,
) {
	data.Label = types.StringValue(pg.Label)
	data.Region = types.StringValue(pg.Region)
	data.PlacementGroupType = types.StringValue(string(pg.PlacementGroupType))
	data.IsCompliant = types.BoolValue(pg.IsCompliant)
	data.PlacementGroupPolicy = types.StringValue(string(pg.PlacementGroupPolicy))

	members := make([]PlacementGroupMemberModel, len(pg.Members))

	for i, member := range pg.Members {
		var m PlacementGroupMemberModel
		m.FlattenMember(member)
		members[i] = m
	}

	data.Members = members

	migrations := pg.Migrations

	if migrations != nil {
		pgMigrations := new(PlacementGroupMigrationModel)
		pgMigrations.FlattenMigrations(*migrations)
		data.Migrations = pgMigrations
	}
}

func (m *PlacementGroupMemberModel) FlattenMember(member linodego.PlacementGroupMember) {
	m.LinodeID = types.Int64Value(int64(member.LinodeID))
	m.IsCompliant = types.BoolValue(member.IsCompliant)
}

func (m *PlacementGroupMigrationModel) FlattenMigrations(migrations linodego.PlacementGroupMigrations) {
	inbound := make([]PlacementGroupMigrationInstanceModel, len(migrations.Inbound))
	outbound := make([]PlacementGroupMigrationInstanceModel, len(migrations.Outbound))

	for i, instance := range migrations.Inbound {
		var m PlacementGroupMigrationInstanceModel
		m.FlattenMigrationInstance(instance)
		inbound[i] = m
	}

	for i, instance := range migrations.Outbound {
		var m PlacementGroupMigrationInstanceModel
		m.FlattenMigrationInstance(instance)
		outbound[i] = m
	}

	m.Inbound = inbound
	m.Outbound = outbound
}

func (m *PlacementGroupMigrationInstanceModel) FlattenMigrationInstance(migrationInstance linodego.PlacementGroupMigrationInstance) {
	m.LinodeID = types.Int64Value(int64(migrationInstance.LinodeID))
}

func (m *PlacementGroupResourceModel) FlattenPlacementGroup(
	ctx context.Context,
	pg *linodego.PlacementGroup,
	preserveKnown bool,
) (resultDiags diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(pg.ID), preserveKnown)

	m.Label = helper.KeepOrUpdateString(m.Label, pg.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateString(m.Region, pg.Region, preserveKnown)
	m.PlacementGroupType = helper.KeepOrUpdateString(m.PlacementGroupType, string(pg.PlacementGroupType), preserveKnown)
	m.PlacementGroupPolicy = helper.KeepOrUpdateString(m.PlacementGroupPolicy, string(pg.PlacementGroupPolicy), preserveKnown)
	m.IsCompliant = helper.KeepOrUpdateBool(m.IsCompliant, pg.IsCompliant, preserveKnown)

	members := make([]PlacementGroupMemberModel, len(pg.Members))
	for i, member := range pg.Members {
		memberModel := PlacementGroupMemberModel{}
		memberModel.FlattenMember(member)
		members[i] = memberModel
	}

	membersSet, diags := types.SetValueFrom(ctx, pgMemberObjectType, members)
	resultDiags.Append(diags...)
	if resultDiags.HasError() {
		return resultDiags
	}

	m.Members = membersSet
	return resultDiags
}

func (m *PlacementGroupResourceModel) CopyFrom(other PlacementGroupResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)

	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Region = helper.KeepOrUpdateValue(m.Region, other.Region, preserveKnown)
	m.PlacementGroupType = helper.KeepOrUpdateValue(m.PlacementGroupType, other.PlacementGroupType, preserveKnown)
	m.PlacementGroupPolicy = helper.KeepOrUpdateValue(m.PlacementGroupPolicy, other.PlacementGroupPolicy, preserveKnown)
	m.IsCompliant = helper.KeepOrUpdateValue(m.IsCompliant, other.IsCompliant, preserveKnown)

	m.Members = other.Members
}
