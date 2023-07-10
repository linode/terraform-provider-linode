package domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

// DomainModel maps a Linode Domain object to a Terraform config.
type DomainModel struct {
	ID          types.Int64                          `tfsdk:"id"`
	Domain      types.String                         `tfsdk:"domain"`
	Type        types.String                         `tfsdk:"type"`
	Group       types.String                         `tfsdk:"group"`
	Status      types.String                         `tfsdk:"status"`
	Description types.String                         `tfsdk:"description"`
	SOAEmail    types.String                         `tfsdk:"soa_email"`
	TTLSec      customtypes.LinodeDomainSecondsValue `tfsdk:"ttl_sec"`
	RetrySec    customtypes.LinodeDomainSecondsValue `tfsdk:"retry_sec"`
	ExpireSec   customtypes.LinodeDomainSecondsValue `tfsdk:"expire_sec"`
	RefreshSec  customtypes.LinodeDomainSecondsValue `tfsdk:"refresh_sec"`
	MasterIPs   types.Set                            `tfsdk:"master_ips"`
	AXFRIPs     types.Set                            `tfsdk:"axfr_ips"`
	Tags        types.Set                            `tfsdk:"tags"`
}

func domainDeepEqual(plan, state DomainModel) bool {
	return state.Domain.Equal(plan.Domain) &&
		state.Type.Equal(plan.Type) &&
		state.Group.Equal(plan.Group) &&
		state.Status.Equal(plan.Status) &&
		state.Description.Equal(plan.Description) &&
		state.SOAEmail.Equal(plan.SOAEmail) &&
		state.TTLSec.Equal(plan.TTLSec) &&
		state.RetrySec.Equal(plan.RetrySec) &&
		state.ExpireSec.Equal(plan.ExpireSec) &&
		state.RefreshSec.Equal(plan.RefreshSec) &&
		state.MasterIPs.Equal(plan.MasterIPs) &&
		state.AXFRIPs.Equal(plan.AXFRIPs) &&
		state.Tags.Equal(plan.Tags)
}

func (m *DomainModel) parseComputed(
	ctx context.Context,
	domain *linodego.Domain,
) diag.Diagnostics {
	m.ID = types.Int64Value(int64(domain.ID))
	m.TTLSec = customtypes.LinodeDomainSecondsValue{
		Int64Value: types.Int64Value(int64(domain.TTLSec)),
	}
	m.RetrySec = customtypes.LinodeDomainSecondsValue{
		Int64Value: types.Int64Value(int64(domain.RetrySec)),
	}
	m.ExpireSec = customtypes.LinodeDomainSecondsValue{
		Int64Value: types.Int64Value(int64(domain.ExpireSec)),
	}
	m.RefreshSec = customtypes.LinodeDomainSecondsValue{
		Int64Value: types.Int64Value(int64(domain.RefreshSec)),
	}

	masterIPs, diags := basetypes.NewSetValueFrom(ctx, types.StringType, domain.MasterIPs)
	if diags.HasError() {
		return diags
	}
	m.MasterIPs = masterIPs

	axfrIPs, diags := basetypes.NewSetValueFrom(ctx, types.StringType, domain.AXfrIPs)
	if diags.HasError() {
		return diags
	}
	m.AXFRIPs = axfrIPs

	var tags basetypes.SetValue
	if len(domain.Tags) == 0 {
		tags = types.SetNull(types.StringType)
	} else {
		tags, diags = basetypes.NewSetValueFrom(ctx, types.StringType, domain.Tags)
		if diags.HasError() {
			return diags
		}
	}
	m.Tags = tags
	return nil
}

func (m *DomainModel) parseNonComputed(
	ctx context.Context,
	domain *linodego.Domain,
) diag.Diagnostics {
	m.Domain = helper.GetValueIfNotNull(domain.Domain)
	m.Type = helper.GetValueIfNotNull(string(domain.Type))
	m.Group = helper.GetValueIfNotNull(domain.Group)
	m.Status = helper.GetValueIfNotNull(string(domain.Status))
	m.Description = helper.GetValueIfNotNull(domain.Description)
	m.SOAEmail = helper.GetValueIfNotNull(domain.SOAEmail)
	return nil
}
