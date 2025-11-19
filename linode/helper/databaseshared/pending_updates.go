package databaseshared

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataSourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ModelPendingUpdate struct {
	Deadline    timetypes.RFC3339 `tfsdk:"deadline"`
	Description types.String      `tfsdk:"description"`
	PlannedFor  timetypes.RFC3339 `tfsdk:"planned_for"`
}

var ResourceAttributePendingUpdates = resourceSchema.SetNestedAttribute{
	Description:   "A set of pending updates.",
	Computed:      true,
	PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	NestedObject: resourceSchema.NestedAttributeObject{
		Attributes: map[string]resourceSchema.Attribute{
			"deadline": resourceSchema.StringAttribute{
				Description: "The time when a mandatory update needs to be applied.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"description": resourceSchema.StringAttribute{
				Description: "A description of the update.",
				Computed:    true,
			},
			"planned_for": resourceSchema.StringAttribute{
				Description: "The date and time a maintenance update will be applied.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	},
}

var DataSourceAttributePendingUpdates = dataSourceSchema.SetNestedAttribute{
	Description: "A set of pending updates.",
	Computed:    true,
	NestedObject: dataSourceSchema.NestedAttributeObject{
		Attributes: map[string]dataSourceSchema.Attribute{
			"deadline": dataSourceSchema.StringAttribute{
				Description: "The time when a mandatory update needs to be applied.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"description": dataSourceSchema.StringAttribute{
				Description: "A description of the update.",
				Computed:    true,
			},
			"planned_for": dataSourceSchema.StringAttribute{
				Description: "The date and time a maintenance update will be applied.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	},
}

var ObjectTypePendingUpdates = ResourceAttributePendingUpdates.NestedObject.Type().(types.ObjectType)

func FlattenPendingUpdates(
	ctx context.Context,
	pending []linodego.DatabaseMaintenanceWindowPending,
) (types.Set, diag.Diagnostics) {
	var d diag.Diagnostics

	pendingObjectsIter := helper.Map(
		slices.Values(pending),
		func(pending linodego.DatabaseMaintenanceWindowPending) attr.Value {
			result, rd := types.ObjectValueFrom(
				ctx,
				ObjectTypePendingUpdates.AttrTypes,
				ModelPendingUpdate{
					Deadline:    timetypes.NewRFC3339TimePointerValue(pending.Deadline),
					Description: types.StringValue(pending.Description),
					PlannedFor:  timetypes.NewRFC3339TimePointerValue(pending.PlannedFor),
				},
			)
			d.Append(rd...)
			return result
		},
	)

	pendingObjectsDeduplicatedIter := helper.FrameworkDropDuplicatesIter(pendingObjectsIter)

	result, rd := types.SetValueFrom(
		ctx,
		ObjectTypePendingUpdates,
		slices.Collect(pendingObjectsDeduplicatedIter),
	)
	d.Append(rd...)
	if d.HasError() {
		return result, d
	}

	return result, nil
}
