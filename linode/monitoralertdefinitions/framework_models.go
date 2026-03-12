package monitoralertdefinitions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinition"
)

// AlertDefinitionFilterModel describes the Terraform resource data model to match the
// resource schema.
type AlertDefinitionFilterModel struct {
	ID               types.String                                            `tfsdk:"id"`
	ServiceType      types.String                                            `tfsdk:"service_type"`
	Filters          frameworkfilter.FiltersModelType                        `tfsdk:"filter"`
	AlertDefinitions []monitoralertdefinition.AlertDefinitionDataSourceModel `tfsdk:"alert_definitions"`
}

func (data *AlertDefinitionFilterModel) parseAlertDefinitions(
	ctx context.Context,
	alertDefinitions []linodego.AlertDefinition,
) diag.Diagnostics {
	result := make([]monitoralertdefinition.AlertDefinitionDataSourceModel, len(alertDefinitions))
	for i := range alertDefinitions {
		var adData monitoralertdefinition.AlertDefinitionDataSourceModel
		diags := adData.FlattenDataSourceModel(ctx, &alertDefinitions[i], false)
		if diags.HasError() {
			return diags
		}
		result[i] = adData
	}

	data.AlertDefinitions = result

	return nil
}
