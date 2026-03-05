package firewallruleset

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// RuleSetBaseModel contains the shared fields used by both the
// resource model and the data-source model.
type RuleSetBaseModel struct {
	Label            types.String         `tfsdk:"label"`
	Description      types.String         `tfsdk:"description"`
	Type             types.String         `tfsdk:"type"`
	Rules            []firewall.RuleModel `tfsdk:"rules"`
	IsServiceDefined types.Bool           `tfsdk:"is_service_defined"`
	Version          types.Int64          `tfsdk:"version"`
	Created          types.String         `tfsdk:"created"`
	Updated          types.String         `tfsdk:"updated"`
}

// FlattenRuleSet maps a linodego.RuleSet into the base model fields.
func (data *RuleSetBaseModel) FlattenRuleSet(
	ctx context.Context,
	rs linodego.RuleSet,
	diags *diag.Diagnostics,
	preserveKnown bool,
) {
	data.Label = helper.KeepOrUpdateString(data.Label, rs.Label, preserveKnown)
	data.Description = helper.KeepOrUpdateString(data.Description, rs.Description, preserveKnown)
	data.Type = helper.KeepOrUpdateString(data.Type, string(rs.Type), preserveKnown)
	data.IsServiceDefined = helper.KeepOrUpdateBool(data.IsServiceDefined, rs.IsServiceDefined, preserveKnown)
	data.Version = helper.KeepOrUpdateInt64(data.Version, int64(rs.Version), preserveKnown)

	if rs.Created != nil {
		data.Created = helper.KeepOrUpdateString(data.Created, rs.Created.Format("2006-01-02T15:04:05"), preserveKnown)
	}
	if rs.Updated != nil {
		data.Updated = helper.KeepOrUpdateString(data.Updated, rs.Updated.Format("2006-01-02T15:04:05"), preserveKnown)
	}

	data.Rules = firewall.FlattenFirewallRules(ctx, rs.Rules, data.Rules, preserveKnown, diags)
}

// RuleSetResourceModel is used by the resource (CRUD).
type RuleSetResourceModel struct {
	ID types.String `tfsdk:"id"`
	RuleSetBaseModel
}

func (data *RuleSetResourceModel) GetCreateOptions(
	ctx context.Context, diags *diag.Diagnostics,
) linodego.RuleSetCreateOptions {
	rules := firewall.ExpandFirewallRules(ctx, data.Rules, diags)

	return linodego.RuleSetCreateOptions{
		Label:       data.Label.ValueString(),
		Description: data.Description.ValueString(),
		Type:        linodego.FirewallRuleSetType(data.Type.ValueString()),
		Rules:       rules,
	}
}

func (data *RuleSetResourceModel) GetUpdateOptions(
	ctx context.Context, diags *diag.Diagnostics,
) linodego.RuleSetUpdateOptions {
	opts := linodego.RuleSetUpdateOptions{}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		v := data.Label.ValueString()
		opts.Label = &v
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		v := data.Description.ValueString()
		opts.Description = &v
	}

	rules := firewall.ExpandFirewallRules(ctx, data.Rules, diags)
	opts.Rules = &rules

	return opts
}

// RuleSetDataSourceModel is used by the single-item data source.
type RuleSetDataSourceModel struct {
	ID types.String `tfsdk:"id"`
	RuleSetBaseModel
}

func (data *RuleSetDataSourceModel) parseRuleSet(
	ctx context.Context,
	rs linodego.RuleSet,
	diags *diag.Diagnostics,
) {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(rs.ID), false)
	data.FlattenRuleSet(ctx, rs, diags, false)
}
