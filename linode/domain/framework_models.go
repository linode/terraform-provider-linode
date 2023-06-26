package domain

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (m *DomainModel) parseDomain(domain *linodego.Domain) {
	m.ID = types.Int64Value(int64(domain.ID))

	m.Domain = types.StringValue(domain.Domain)
	m.Type = types.StringValue(string(domain.Type))
	m.Group = types.StringValue(domain.Group)
	m.Status = types.StringValue(string(domain.Status))
	m.Description = types.StringValue(domain.Description)
	m.SOAEmail = types.StringValue(domain.SOAEmail)

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

	m.MasterIPs = helper.StringSliceToFrameworkSet(domain.MasterIPs)
	m.AXFRIPs = helper.StringSliceToFrameworkSet(domain.AXfrIPs)
	m.Tags = helper.StringSliceToFrameworkSet(domain.Tags)
}
