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

func (m *DomainModel) parseDomain(
	ctx context.Context,
	domain *linodego.Domain,
) diag.Diagnostics {
	m.ID = types.Int64Value(int64(domain.ID))

	m.Domain = helper.GetValueIfNotNull(domain.Domain)
	m.Type = helper.GetValueIfNotNull(string(domain.Type))
	m.Group = helper.GetValueIfNotNull(domain.Group)
	m.Status = helper.GetValueIfNotNull(string(domain.Status))
	m.Description = helper.GetValueIfNotNull(domain.Description)
	m.SOAEmail = helper.GetValueIfNotNull(domain.SOAEmail)

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

	tags, diags := basetypes.NewSetValueFrom(ctx, types.StringType, domain.Tags)
	if diags.HasError() {
		return diags
	}
	m.Tags = tags

	return nil
}
