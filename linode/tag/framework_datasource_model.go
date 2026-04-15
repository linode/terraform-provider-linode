package tag

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Label   types.String `tfsdk:"label"`
	Objects types.List   `tfsdk:"objects"`
}

type TaggedObjectModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
}

func (data *DataSourceModel) FlattenTaggedObjects(
	ctx context.Context,
	objects linodego.TaggedObjectList,
	diags *diag.Diagnostics,
) {
	models := make([]TaggedObjectModel, 0, len(objects))

	for _, obj := range objects {
		m := TaggedObjectModel{
			Type: types.StringValue(obj.Type),
		}

		switch obj.Type {
		case "linode":
			if inst, ok := obj.Data.(linodego.Instance); ok {
				m.ID = types.StringValue(strconv.Itoa(inst.ID))
			}
		case "domain":
			if d, ok := obj.Data.(linodego.Domain); ok {
				m.ID = types.StringValue(strconv.Itoa(d.ID))
			}
		case "volume":
			if v, ok := obj.Data.(linodego.Volume); ok {
				m.ID = types.StringValue(strconv.Itoa(v.ID))
			}
		case "nodebalancer":
			if n, ok := obj.Data.(linodego.NodeBalancer); ok {
				m.ID = types.StringValue(strconv.Itoa(n.ID))
			}
		case "reserved_ipv4_address":
			if ip, ok := obj.Data.(linodego.InstanceIP); ok {
				m.ID = types.StringValue(ip.Address)
			}
		default:
			diags.AddWarning("Unknown tagged object type",
				fmt.Sprintf("tagged object type %q is not recognised; ID will be empty", obj.Type))
			m.ID = types.StringValue("")
		}

		models = append(models, m)
	}

	listVal, d := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: tagObjectAttrTypes},
		models,
	)
	diags.Append(d...)
	if !d.HasError() {
		data.Objects = listVal
	}
}
