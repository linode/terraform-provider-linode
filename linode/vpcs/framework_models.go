package vpcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/vpc"
)

// VPCFilterModel describes the Terraform resource data model to match the
// resource schema.
type VPCFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCs    []vpc.VPCModel                   `tfsdk:"vpcs"`
}

func (model *VPCFilterModel) FlattenVPCs(
	ctx context.Context,
	vpcs []linodego.VPC,
	preserveKnown bool,
) {
	vpcModels := make([]vpc.VPCModel, len(vpcs))

	for i := range vpcs {
		var vpc vpc.VPCModel
		vpc.FlattenVPC(ctx, &vpcs[i], preserveKnown)
		vpcModels[i] = vpc
	}

	model.VPCs = vpcModels
}
