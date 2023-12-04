package vlan

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type VLANModel struct {
	Label   types.String `tfsdk:"label"`
	Linodes types.Set    `tfsdk:"linodes"`
	Region  types.String `tfsdk:"region"`
	Created types.String `tfsdk:"created"`
}

type VLANsFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	VLANs   []VLANModel                      `tfsdk:"vlans"`
}

func (data *VLANsFilterModel) parseVLANs(
	ctx context.Context,
	vlans []linodego.VLAN,
) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := json.Marshal(vlans)
	if err != nil {
		diags.AddError("Error marshalling json", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))
	results := make([]VLANModel, len(vlans))

	for i, v := range vlans {
		var vlan VLANModel
		vlan.parseVLAN(ctx, v)
		results[i] = vlan
	}

	data.VLANs = results

	return diags
}

func (data *VLANModel) parseVLAN(
	ctx context.Context,
	vlan linodego.VLAN,
) diag.Diagnostics {
	data.Label = types.StringValue(vlan.Label)

	linodes, diags := types.SetValueFrom(ctx, types.Int64Type, vlan.Linodes)
	if diags.HasError() {
		return diags
	}

	data.Linodes = linodes
	data.Region = types.StringValue(vlan.Region)

	if vlan.Created != nil {
		data.Created = types.StringValue(vlan.Created.Format(time.RFC3339))
	}

	return diags
}
