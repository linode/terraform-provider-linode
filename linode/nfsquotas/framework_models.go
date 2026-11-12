package nfsquotas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsquota"
)

// NFSQuotaFilterModel describes the Terraform data model to match the
// data source schema.
type NFSQuotaFilterModel struct {
	ID           types.String                     `tfsdk:"id"`
	SpaceID      types.Int64                      `tfsdk:"space_id"`
	FilesystemID types.Int64                      `tfsdk:"filesystem_id"`
	Filters      frameworkfilter.FiltersModelType `tfsdk:"filter"`
	NFSQuotas    []nfsquota.DataSourceModel       `tfsdk:"nfs_quotas"`
}

func (model *NFSQuotaFilterModel) FlattenNFSQuotas(quotas []linodego.NFSQuota) {
	result := make([]nfsquota.DataSourceModel, len(quotas))

	for i := range quotas {
		var m nfsquota.DataSourceModel
		m.ParseQuota(&quotas[i])
		m.SpaceID = model.SpaceID
		result[i] = m
	}

	model.NFSQuotas = result
}
