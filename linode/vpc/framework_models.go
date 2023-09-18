package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

type VPCModel struct {
	ID          types.Int64                        `tfsdk:"id"`
	Label       types.String                       `tfsdk:"label"`
	Description types.String                       `tfsdk:"description"`
	Region      types.String                       `tfsdk:"region"`
	Created     customtypes.RFC3339TimeStringValue `tfsdk:"created"`
	Updated     customtypes.RFC3339TimeStringValue `tfsdk:"updated"`
}

func (d *VPCModel) parseComputedAttributes(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(vpc.ID))
	d.Description = types.StringValue(vpc.Description)
	d.Created = customtypes.RFC3339TimeStringValue{
		StringValue: helper.NullableTimeToFramework(vpc.Created),
	}
	d.Updated = customtypes.RFC3339TimeStringValue{
		StringValue: helper.NullableTimeToFramework(vpc.Updated),
	}

	return nil
}

func (d *VPCModel) parseVPC(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.Label = types.StringValue(vpc.Label)
	d.Region = types.StringValue(vpc.Region)

	return d.parseComputedAttributes(ctx, vpc)
}
