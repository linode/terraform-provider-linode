package nfsfilesystems

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsfilesystem"
)

// NFSFilesystemFilterModel describes the Terraform resource data model to match the
// resource schema.
type NFSFilesystemFilterModel struct {
	ID             types.String                     `tfsdk:"id"`
	SpaceID        types.Int64                      `tfsdk:"space_id"`
	Filters        frameworkfilter.FiltersModelType `tfsdk:"filter"`
	NFSFilesystems []nfsfilesystem.DataSourceModel  `tfsdk:"nfs_filesystems"`
}

func (model *NFSFilesystemFilterModel) FlattenNFSFilesystems(
	ctx context.Context,
	filesystems []linodego.NFSFilesystem,
	preserveKnown bool,
) diag.Diagnostics {
	result := make([]nfsfilesystem.DataSourceModel, len(filesystems))

	for i := range filesystems {
		var m nfsfilesystem.DataSourceModel
		if d := m.Flatten(ctx, &filesystems[i], preserveKnown); d.HasError() {
			return d
		}
		result[i] = m
	}

	model.NFSFilesystems = result
	return nil
}
