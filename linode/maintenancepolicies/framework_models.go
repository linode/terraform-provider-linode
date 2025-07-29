package maintenancepolicies

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type MaintenancePolicyModel struct {
	Slug                  types.String `tfsdk:"slug"`
	Label                 types.String `tfsdk:"label"`
	Description           types.String `tfsdk:"description"`
	Type                  types.String `tfsdk:"type"`
	NotificationPeriodSec types.Int64  `tfsdk:"notification_period_sec"`
	IsDefault             types.Bool   `tfsdk:"is_default"`
}

type MaintenancePolicyFilterModel struct {
	ID                  types.String                     `tfsdk:"id"`
	Filters             frameworkfilter.FiltersModelType `tfsdk:"filter"`
	MaintenancePolicies []MaintenancePolicyModel         `tfsdk:"maintenance_policies"`
}

func (model *MaintenancePolicyFilterModel) parseMaintenancePolicies(maintenancePolicies []linodego.MaintenancePolicy) diag.Diagnostics {
	result := make([]MaintenancePolicyModel, len(maintenancePolicies))

	for i := range maintenancePolicies {
		var m MaintenancePolicyModel

		diags := m.parseMaintenancePolicy(&maintenancePolicies[i])
		if diags.HasError() {
			return diags
		}

		result[i] = m
	}

	model.MaintenancePolicies = result

	id, _ := json.Marshal(maintenancePolicies)
	model.ID = types.StringValue(string(id))

	return nil
}

func (data *MaintenancePolicyModel) parseMaintenancePolicy(maintenancePolicy *linodego.MaintenancePolicy,
) diag.Diagnostics {
	data.Slug = types.StringValue(maintenancePolicy.Slug)
	data.Label = types.StringValue(maintenancePolicy.Label)
	data.Description = types.StringValue(maintenancePolicy.Description)
	data.Type = types.StringValue(maintenancePolicy.Type)
	data.NotificationPeriodSec = types.Int64Value(int64(maintenancePolicy.NotificationPeriodSec))
	data.IsDefault = types.BoolValue(maintenancePolicy.IsDefault)

	return nil
}
