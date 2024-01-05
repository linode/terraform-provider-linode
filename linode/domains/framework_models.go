package domains

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/domain"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

// DomainFilterModel describes the Terraform resource data model to match the
// resource schema.
type DomainFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Domains []domain.DomainModel             `tfsdk:"domains"`
}

func (data *DomainFilterModel) parseDomains(
	domains []linodego.Domain,
) {
	result := make([]domain.DomainModel, len(domains))
	for i := range domains {
		var domainData domain.DomainModel
		domainData.ParseDomain(&domains[i])
		result[i] = domainData
	}

	data.Domains = result
}
