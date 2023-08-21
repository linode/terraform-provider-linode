package sshkeys

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/sshkey"
)

// NodeBalancerFilterModel describes the Terraform resource data model to match the
// resource schema.
type SSHKeyFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	SSHKeys []sshkey.DataSourceModel         `tfsdk:"sshkeys"`
}

func (data *SSHKeyFilterModel) parseSSHKeys(
	ctx context.Context,
	sshkeys []linodego.SSHKey,
) {
	result := make([]sshkey.DataSourceModel, len(sshkeys))
	for i := range sshkeys {
		var sshData sshkey.DataSourceModel
		sshData.ParseSSHKey(&sshkeys[i])
		result[i] = sshData
	}

	data.SSHKeys = result
}
