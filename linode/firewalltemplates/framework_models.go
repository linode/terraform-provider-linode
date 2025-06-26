package firewalltemplates

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalltemplate"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/frameworkfilter"
)

type FirewallTemplateFilterModel struct {
	ID                types.String                                 `tfsdk:"id"`
	Filters           frameworkfilter.FiltersModelType             `tfsdk:"filter"`
	FirewallTemplates []firewalltemplate.FirewallTemplateBaseModel `tfsdk:"firewall_templates"`
}

func (data *FirewallTemplateFilterModel) parseFirewallTemplates(
	ctx context.Context,
	templates []linodego.FirewallTemplate,
	diags *diag.Diagnostics,
) {
	data.FirewallTemplates = make([]firewalltemplate.FirewallTemplateBaseModel, len(templates))
	for i, t := range templates {
		data.FirewallTemplates[i].FlattenFirewallTemplate(ctx, t, diags, false)
	}
}
