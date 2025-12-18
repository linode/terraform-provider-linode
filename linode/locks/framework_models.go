package locks

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

// LockModel describes a single Lock.
type LockModel struct {
	ID          types.Int64  `tfsdk:"id"`
	EntityID    types.Int64  `tfsdk:"entity_id"`
	EntityType  types.String `tfsdk:"entity_type"`
	LockType    types.String `tfsdk:"lock_type"`
	EntityLabel types.String `tfsdk:"entity_label"`
	EntityURL   types.String `tfsdk:"entity_url"`
}

func (m *LockModel) parseLock(lock linodego.Lock) {
	m.ID = types.Int64Value(int64(lock.ID))
	m.LockType = types.StringValue(string(lock.LockType))
	m.EntityID = types.Int64Value(int64(lock.Entity.ID))
	m.EntityType = types.StringValue(string(lock.Entity.Type))
	m.EntityLabel = types.StringValue(lock.Entity.Label)
	m.EntityURL = types.StringValue(lock.Entity.URL)
}

// LockFilterModel describes the Terraform data source data model to match the
// resource schema.
type LockFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Locks   []LockModel                      `tfsdk:"locks"`
}

func (data *LockFilterModel) parseLocks(locks []linodego.Lock) diag.Diagnostics {
	result := make([]LockModel, len(locks))

	for i := range locks {
		var lockData LockModel
		lockData.parseLock(locks[i])
		result[i] = lockData
	}

	data.Locks = result
	return nil
}
