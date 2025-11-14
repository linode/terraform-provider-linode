package databaseshared

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	dataSourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type ModelUpdates struct {
	DayOfWeek types.Int64  `tfsdk:"day_of_week"`
	Duration  types.Int64  `tfsdk:"duration"`
	Frequency types.String `tfsdk:"frequency"`
	HourOfDay types.Int64  `tfsdk:"hour_of_day"`
}

func (m ModelUpdates) ToLinodego(d diag.Diagnostics) *linodego.DatabaseMaintenanceWindow {
	return &linodego.DatabaseMaintenanceWindow{
		DayOfWeek: linodego.DatabaseDayOfWeek(helper.FrameworkSafeInt64ToInt(m.DayOfWeek.ValueInt64(), &d)),
		Duration:  helper.FrameworkSafeInt64ToInt(m.Duration.ValueInt64(), &d),
		Frequency: linodego.DatabaseMaintenanceFrequency(m.Frequency.ValueString()),
		HourOfDay: helper.FrameworkSafeInt64ToInt(m.HourOfDay.ValueInt64(), &d),
	}
}

var ResourceAttributeUpdates = resourceSchema.SingleNestedAttribute{
	Description: "Configuration settings for automated patch update maintenance for the Managed Database.",
	Attributes: map[string]resourceSchema.Attribute{
		"day_of_week": resourceSchema.Int64Attribute{
			Description: "The numeric reference for the day of the week to perform maintenance. " +
				"1 is Monday, 2 is Tuesday, through to 7 which is Sunday.",
			Optional: true,
			Computed: true,
			Validators: []validator.Int64{
				int64validator.Between(1, 7),
			},
		},
		"duration": resourceSchema.Int64Attribute{
			Description: "The maximum maintenance window time in hours.",
			Optional:    true,
			Computed:    true,
		},
		"frequency": resourceSchema.StringAttribute{
			Description: "How frequently maintenance occurs. Currently can only be weekly.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("weekly"),
		},
		"hour_of_day": resourceSchema.Int64Attribute{
			Description: "How frequently maintenance occurs. Currently can only be weekly.",
			Optional:    true,
			Computed:    true,
			Validators: []validator.Int64{
				int64validator.Between(0, 23),
			},
		},
	},
	Computed:      true,
	Optional:      true,
	PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
}

var DataSourceAttributeUpdates = dataSourceSchema.SingleNestedAttribute{
	Description: "Configuration settings for automated patch update maintenance for the Managed Database.",
	Attributes: map[string]dataSourceSchema.Attribute{
		"day_of_week": dataSourceSchema.Int64Attribute{
			Description: "The numeric reference for the day of the week to perform maintenance. " +
				"1 is Monday, 2 is Tuesday, through to 7 which is Sunday.",
			Computed: true,
		},
		"duration": dataSourceSchema.Int64Attribute{
			Description: "The maximum maintenance window time in hours.",
			Computed:    true,
		},
		"frequency": dataSourceSchema.StringAttribute{
			Description: "How frequently maintenance occurs. Currently can only be weekly.",
			Computed:    true,
		},
		"hour_of_day": dataSourceSchema.Int64Attribute{
			Description: "How frequently maintenance occurs. Currently can only be weekly.",
			Computed:    true,
			Validators: []validator.Int64{
				int64validator.Between(0, 23),
			},
		},
	},
	Computed: true,
}

var ObjectTypeUpdates = ResourceAttributeUpdates.GetType().(types.ObjectType)

func FlattenUpdates(ctx context.Context, updates linodego.DatabaseMaintenanceWindow) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(
		ctx,
		ObjectTypeUpdates.AttrTypes,
		&ModelUpdates{
			DayOfWeek: types.Int64Value(int64(updates.DayOfWeek)),
			Duration:  types.Int64Value(int64(updates.Duration)),
			Frequency: types.StringValue(string(updates.Frequency)),
			HourOfDay: types.Int64Value(int64(updates.HourOfDay)),
		},
	)
}
