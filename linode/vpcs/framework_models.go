package vpcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/vpc"
)

// VPCFilterModel describes the Terraform resource data model to match the
// resource schema.
type VPCFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	VPCs    []vpc.Model                      `tfsdk:"vpcs"`
}

func (model *VPCFilterModel) FlattenVPCs(
	ctx context.Context,
	vpcs []linodego.VPC,
	preserveKnown bool,
) diag.Diagnostics {
	vpcModels := make([]vpc.Model, len(vpcs))

	for i := range vpcs {
		var vpc vpc.Model
		if d := vpc.FlattenVPC(ctx, &vpcs[i], preserveKnown); d.HasError() {
			return d
		}
		vpcModels[i] = vpc
	}

	model.VPCs = vpcModels
	return nil
}
