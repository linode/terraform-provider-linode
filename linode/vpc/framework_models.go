package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type VPCModel struct {
	ID          types.Int64       `tfsdk:"id"`
	Label       types.String      `tfsdk:"label"`
	Description types.String      `tfsdk:"description"`
	Region      types.String      `tfsdk:"region"`
	Created     timetypes.RFC3339 `tfsdk:"created"`
	Updated     timetypes.RFC3339 `tfsdk:"updated"`
}

func (d *VPCModel) parseComputedAttributes(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.ID = types.Int64Value(int64(vpc.ID))
	d.Description = types.StringValue(vpc.Description)
	d.Created = timetypes.NewRFC3339TimePointerValue(vpc.Created)
	d.Updated = timetypes.NewRFC3339TimePointerValue(vpc.Updated)
	return nil
}

func (d *VPCModel) ParseVPC(
	ctx context.Context,
	vpc *linodego.VPC,
) diag.Diagnostics {
	d.Label = types.StringValue(vpc.Label)
	d.Region = types.StringValue(vpc.Region)

	return d.parseComputedAttributes(ctx, vpc)
}
