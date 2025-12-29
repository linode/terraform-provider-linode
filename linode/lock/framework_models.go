package lock

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// ResourceModel describes the Terraform resource model to match the
// resource schema.
type ResourceModel struct {
	ID         types.String `tfsdk:"id"`
	EntityID   types.Int64  `tfsdk:"entity_id"`
	EntityType types.String `tfsdk:"entity_type"`
	LockType   types.String `tfsdk:"lock_type"`

	// Computed fields from the Entity
	EntityLabel types.String `tfsdk:"entity_label"`
	EntityURL   types.String `tfsdk:"entity_url"`
}

func (m *ResourceModel) FlattenLock(lock *linodego.Lock, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(lock.ID), preserveKnown)
	m.LockType = helper.KeepOrUpdateString(m.LockType, string(lock.LockType), preserveKnown)

	// Flatten the entity information
	m.EntityID = helper.KeepOrUpdateInt64(m.EntityID, int64(lock.Entity.ID), preserveKnown)
	m.EntityType = helper.KeepOrUpdateString(m.EntityType, string(lock.Entity.Type), preserveKnown)
	m.EntityLabel = helper.KeepOrUpdateString(m.EntityLabel, lock.Entity.Label, preserveKnown)
	m.EntityURL = helper.KeepOrUpdateString(m.EntityURL, lock.Entity.URL, preserveKnown)
}

func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.EntityID = helper.KeepOrUpdateValue(m.EntityID, other.EntityID, preserveKnown)
	m.EntityType = helper.KeepOrUpdateValue(m.EntityType, other.EntityType, preserveKnown)
	m.LockType = helper.KeepOrUpdateValue(m.LockType, other.LockType, preserveKnown)
	m.EntityLabel = helper.KeepOrUpdateValue(m.EntityLabel, other.EntityLabel, preserveKnown)
	m.EntityURL = helper.KeepOrUpdateValue(m.EntityURL, other.EntityURL, preserveKnown)
}

func (m *ResourceModel) GetCreateOptions(diags *diag.Diagnostics) linodego.LockCreateOptions {
	entityID := helper.FrameworkSafeInt64ToInt(m.EntityID.ValueInt64(), diags)
	if diags.HasError() {
		return linodego.LockCreateOptions{}
	}

	return linodego.LockCreateOptions{
		EntityID:   entityID,
		EntityType: linodego.EntityType(m.EntityType.ValueString()),
		LockType:   linodego.LockType(m.LockType.ValueString()),
	}
}
