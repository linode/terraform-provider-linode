package nfsspaces

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsspace"
)

// NFSSpaceFilterModel describes the Terraform resource data model to match the
// resource schema.
type NFSSpaceFilterModel struct {
	ID        types.String                     `tfsdk:"id"`
	Filters   frameworkfilter.FiltersModelType `tfsdk:"filter"`
	NFSSpaces []nfsspace.DataSourceModel       `tfsdk:"nfs_spaces"`
}

func (model *NFSSpaceFilterModel) FlattenNFSSpaces(
	ctx context.Context,
	spaces []linodego.NFSSpace,
	preserveKnown bool,
) diag.Diagnostics {
	result := make([]nfsspace.DataSourceModel, len(spaces))

	for i := range spaces {
		var m nfsspace.DataSourceModel
		if d := m.Flatten(ctx, &spaces[i], preserveKnown); d.HasError() {
			return d
		}
		result[i] = m
	}

	model.NFSSpaces = result
	return nil
}
