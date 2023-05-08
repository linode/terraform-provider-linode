package domain

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

// DomainModel describes the Terraform resource data model to match the
// resource schema.
type DomainModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain" linode_mutable:"true"`
	Type        types.String `tfsdk:"type" linode_mutable:"true"`
	Group       types.String `tfsdk:"group" linode_mutable:"true"`
	Status      types.String `tfsdk:"status" linode_mutable:"true"`
	Description types.String `tfsdk:"description" linode_mutable:"true"`
	MasterIPs   types.Set    `tfsdk:"master_ips" linode_mutable:"true"`
	AXFRIPs     types.Set    `tfsdk:"axfr_ips" linode_mutable:"true"`
	TTLSec      types.Int64  `tfsdk:"ttl_sec" linode_mutable:"true"`
	RetrySec    types.Int64  `tfsdk:"retry_sec" linode_mutable:"true"`
	ExpireSec   types.Int64  `tfsdk:"expire_sec" linode_mutable:"true"`
	RefreshSec  types.Int64  `tfsdk:"refresh_sec" linode_mutable:"true"`
	SOAEmail    types.String `tfsdk:"soa_email" linode_mutable:"true"`
	Tags        types.Set    `tfsdk:"tags" linode_mutable:"true"`
}

func (data *DomainModel) parseDomain(
	ctx context.Context,
	domain *linodego.Domain,
) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(domain.ID))
	data.Domain = types.StringValue(domain.Domain)
	data.Type = types.StringValue(string(domain.Type))
	data.Group = types.StringValue(domain.Group)
	data.Status = types.StringValue(string(domain.Status))
	data.Description = types.StringValue(domain.Description)
	data.TTLSec = types.Int64Value(int64(domain.TTLSec))
	data.RetrySec = types.Int64Value(int64(domain.RetrySec))
	data.ExpireSec = types.Int64Value(int64(domain.ExpireSec))
	data.RefreshSec = types.Int64Value(int64(domain.RefreshSec))
	data.SOAEmail = types.StringValue(domain.SOAEmail)

	masterIPs, err := types.SetValueFrom(ctx, types.StringType, domain.MasterIPs)
	diagnostics.Append(err...)
	if diagnostics.HasError() {
		return diagnostics
	}
	data.MasterIPs = masterIPs

	axfrIPs, err := types.SetValueFrom(ctx, types.StringType, domain.AXfrIPs)
	diagnostics.Append(err...)
	if diagnostics.HasError() {
		return diagnostics
	}
	data.AXFRIPs = axfrIPs

	tags, err := types.SetValueFrom(ctx, types.StringType, domain.Tags)
	diagnostics.Append(err...)
	if diagnostics.HasError() {
		return diagnostics
	}
	data.Tags = tags

	return nil
}
