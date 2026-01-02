package lock

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// DataSourceModel describes the Terraform data source model.
type DataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	EntityID    types.Int64  `tfsdk:"entity_id"`
	EntityType  types.String `tfsdk:"entity_type"`
	LockType    types.String `tfsdk:"lock_type"`
	EntityLabel types.String `tfsdk:"entity_label"`
	EntityURL   types.String `tfsdk:"entity_url"`
}

func (m *DataSourceModel) ParseLock(lock *linodego.Lock) {
	m.ID = types.Int64Value(int64(lock.ID))
	m.LockType = helper.KeepOrUpdateString(m.LockType, string(lock.LockType), false)
	m.EntityID = helper.KeepOrUpdateInt64(m.EntityID, int64(lock.Entity.ID), false)
	m.EntityType = helper.KeepOrUpdateString(m.EntityType, string(lock.Entity.Type), false)
	m.EntityLabel = helper.KeepOrUpdateString(m.EntityLabel, lock.Entity.Label, false)
	m.EntityURL = helper.KeepOrUpdateString(m.EntityURL, lock.Entity.URL, false)
}
