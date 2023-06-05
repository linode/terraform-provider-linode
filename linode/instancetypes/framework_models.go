package instancetypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
)

type InstanceTypeModel struct {
	ID         types.String `tfsdk:"id"`
	Label      types.String `tfsdk:"label"`
	Disk       types.Int64  `tfsdk:"disk"`
	Class      types.String `tfsdk:"class"`
	Price      types.List   `tfsdk:"price"`
	Addons     types.List   `tfsdk:"addons"`
	NetworkOut types.Int64  `tfsdk:"network_out"`
	Memory     types.Int64  `tfsdk:"memory"`
	Transfer   types.Int64  `tfsdk:"transfer"`
	VCPUs      types.Int64  `tfsdk:"vcpus"`
}

type InstanceTypeFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Types   []InstanceTypeModel              `tfsdk:"types"`
}

func (model *InstanceTypeFilterModel) parseInstanceTypes(ctx context.Context, instanceTypes []linodego.LinodeType) diag.Diagnostics {
	parseInstanceType := func(instanceType linodego.LinodeType) (InstanceTypeModel, diag.Diagnostics) {
		var m InstanceTypeModel

		// Gonna need to add more parse functions for the object types

		// m.Datetime = types.StringValue(login.Datetime.Format(time.RFC3339))
		// m.ID = types.Int64Value(int64(login.ID))
		// m.IP = types.StringValue(login.IP)
		// m.Restricted = types.BoolValue(login.Restricted)
		// m.Username = types.StringValue(login.Username)
		// m.Status = types.StringValue(login.Status)
		m.ID = types.StringValue(instanceType.ID)
		m.Disk = types.Int64Value(int64(instanceType.Disk))
		m.Class = types.StringValue(string(instanceType.Class))

		price, diags := instancetype.FlattenPrice(ctx, *instanceType.Price)
		if diags.HasError() {
			return InstanceTypeModel{}, diags
		}
		m.Price = *price

		m.Label = types.StringValue(instanceType.Label)

		addons, diags := instancetype.FlattenAddons(ctx, *instanceType.Addons)
		if diags.HasError() {
			return InstanceTypeModel{}, diags
		}
		m.Addons = *addons

		m.NetworkOut = types.Int64Value(int64(instanceType.NetworkOut))
		m.Memory = types.Int64Value(int64(instanceType.Memory))
		m.Transfer = types.Int64Value(int64(instanceType.Transfer))
		m.VCPUs = types.Int64Value(int64(instanceType.VCPUs))

		return m, nil
	}

	result := make([]InstanceTypeModel, len(instanceTypes))

	for i, instanceType := range instanceTypes {

		res, diags := parseInstanceType(instanceType)
		if diags.HasError() {
			return diags
		}

		result[i] = res
	}

	model.Types = result

	return nil
}
