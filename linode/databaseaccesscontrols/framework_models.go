package databaseaccesscontrols

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// ResourceModel represents the Terraform state of a
// linode_database_access_controls resource.
type ResourceModel struct {
	ID           types.String `tfsdk:"id"`
	DatabaseID   types.Int64  `tfsdk:"database_id"`
	DatabaseType types.String `tfsdk:"database_type"`
	AllowList    types.Set    `tfsdk:"allow_list"`
}

// Flatten populates this model from the given database identity and allow list.
func (m *ResourceModel) Flatten(
	dbID int,
	dbType string,
	allowList []string,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	m.ID = helper.KeepOrUpdateString(m.ID, formatID(dbID, dbType), preserveKnown)
	m.DatabaseID = helper.KeepOrUpdateInt64(m.DatabaseID, int64(dbID), preserveKnown)
	m.DatabaseType = helper.KeepOrUpdateString(m.DatabaseType, dbType, preserveKnown)
	m.AllowList = helper.KeepOrUpdateStringSet(m.AllowList, allowList, preserveKnown, diags)
}

// CopyFrom copies the values of the given model into this model.
func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.DatabaseID = helper.KeepOrUpdateValue(m.DatabaseID, other.DatabaseID, preserveKnown)
	m.DatabaseType = helper.KeepOrUpdateValue(m.DatabaseType, other.DatabaseType, preserveKnown)
	m.AllowList = helper.KeepOrUpdateValue(m.AllowList, other.AllowList, preserveKnown)
}
