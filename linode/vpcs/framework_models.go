package vpcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/vpc"
)

// VPCFilterModel describes the Terraform resource data model to match the
// resource schema.
type VPCFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	VPCs    []vpc.VPCModel                   `tfsdk:"vpcs"`
}

func (model *VPCFilterModel) parseVPCs(
	ctx context.Context,
	vpcs []linodego.VPC,
) diag.Diagnostics {
	vpcModels := make([]vpc.VPCModel, len(vpcs))

	for i := range vpcs {
		var vpc vpc.VPCModel

		diags := vpc.ParseVPC(ctx, &vpcs[i])
		if diags.HasError() {
			return diags
		}

		vpcModels[i] = vpc
	}

	model.VPCs = vpcModels

	return nil
}
