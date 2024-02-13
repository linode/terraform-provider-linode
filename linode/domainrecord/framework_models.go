package domainrecord

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"
)

type DataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	DomainID types.Int64  `tfsdk:"domain_id"`
	Type     types.String `tfsdk:"type"`
	TTLSec   types.Int64  `tfsdk:"ttl_sec"`
	Target   types.String `tfsdk:"target"`
	Priority types.Int64  `tfsdk:"priority"`
	Weight   types.Int64  `tfsdk:"weight"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Service  types.String `tfsdk:"service"`
	Tag      types.String `tfsdk:"tag"`
}

func (data *DataSourceModel) FlattenDomainRecord(domainRecord *linodego.DomainRecord) {
	data.ID = types.Int64Value(int64(domainRecord.ID))
	data.Name = types.StringValue(domainRecord.Name)
	data.Type = types.StringValue(string(domainRecord.Type))
	data.TTLSec = types.Int64Value(int64(domainRecord.TTLSec))
	data.Target = types.StringValue(domainRecord.Target)
	data.Priority = types.Int64Value(int64(domainRecord.Priority))
	data.Weight = types.Int64Value(int64(domainRecord.Weight))
	data.Port = types.Int64Value(int64(domainRecord.Port))
	data.Protocol = types.StringPointerValue(domainRecord.Protocol)
	data.Service = types.StringPointerValue(domainRecord.Service)
	data.Tag = types.StringPointerValue(domainRecord.Tag)
}

type ResourceModel struct {
	ID         types.Int64                         `tfsdk:"id"`
	Name       types.String                        `tfsdk:"name"`
	DomainID   types.Int64                         `tfsdk:"domain_id"`
	RecordType types.String                        `tfsdk:"record_type"`
	TTLSec     customtypes.DomainRecordTTLValue    `tfsdk:"ttl_sec"`
	Target     customtypes.DomainRecordTargetValue `tfsdk:"target"`
	Priority   types.Int64                         `tfsdk:"priority"`
	Weight     types.Int64                         `tfsdk:"weight"`
	Port       types.Int64                         `tfsdk:"port"`
	Protocol   types.String                        `tfsdk:"protocol"`
	Service    types.String                        `tfsdk:"service"`
	Tag        types.String                        `tfsdk:"tag"`
}

func (data *ResourceModel) DomainRecordNameSemanticEquals(
	ctx context.Context,
	client *linodego.Client,
	domainID int,
	oldValue, newValue string,
	diags *diag.Diagnostics,
) bool {
	tflog.Trace(ctx, "client.GetDomain(...)")
	domain, err := client.GetDomain(ctx, domainID)
	if err != nil {
		diags.AddError("Failed to get parent domain", err.Error())
		return false
	}

	return strings.TrimSuffix(strings.TrimSuffix(oldValue, domain.Domain), ".") == newValue
}

func (data *ResourceModel) FlattenDomainRecord(
	ctx context.Context,
	client *linodego.Client,
	domainRecord *linodego.DomainRecord,
	preserveKnown bool,
) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = helper.KeepOrUpdateInt64(data.ID, int64(domainRecord.ID), preserveKnown)

	domainID := helper.FrameworkSafeInt64ToInt(data.DomainID.ValueInt64(), &diags)
	if diags.HasError() {
		return diags
	}

	// Only to update name when it's null, not configured, or semantic different than the new value.
	// If the user has specified their FQDN as a part of their planned record name but the API trims
	// the FQDN from the returned record name, then they are semantic equals.
	if data.Name.IsNull() || data.Name.IsUnknown() || !data.DomainRecordNameSemanticEquals(
		ctx, client, domainID, data.Name.ValueString(), domainRecord.Name, &diags,
	) {
		data.Name = helper.KeepOrUpdateString(data.Name, domainRecord.Name, preserveKnown)
	}
	if diags.HasError() {
		return diags
	}

	data.RecordType = helper.KeepOrUpdateString(data.RecordType, string(domainRecord.Type), preserveKnown)

	data.TTLSec = helper.KeepOrUpdateValue(
		data.TTLSec,
		customtypes.DomainRecordTTLValue{
			Int64Value: types.Int64Value(int64(domainRecord.TTLSec)),
		},
		preserveKnown,
	)

	data.Target = helper.KeepOrUpdateValue(
		data.Target,
		customtypes.DomainRecordTargetValue{
			StringValue: types.StringValue(domainRecord.Target),
		},
		preserveKnown,
	)

	data.Priority = helper.KeepOrUpdateInt64(data.Priority, int64(domainRecord.Priority), preserveKnown)
	data.Weight = helper.KeepOrUpdateInt64(data.Weight, int64(domainRecord.Weight), preserveKnown)
	data.Port = helper.KeepOrUpdateInt64(data.Port, int64(domainRecord.Port), preserveKnown)
	data.Protocol = helper.KeepOrUpdateStringPointer(data.Protocol, domainRecord.Protocol, preserveKnown)
	data.Service = helper.KeepOrUpdateStringPointer(data.Service, domainRecord.Service, preserveKnown)
	data.Tag = helper.KeepOrUpdateStringPointer(data.Tag, domainRecord.Tag, preserveKnown)

	return diags
}

func (data *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.Name = helper.KeepOrUpdateValue(data.Name, other.Name, preserveKnown)
	data.DomainID = helper.KeepOrUpdateValue(data.DomainID, other.DomainID, preserveKnown)
	data.RecordType = helper.KeepOrUpdateValue(data.RecordType, other.RecordType, preserveKnown)
	data.TTLSec = helper.KeepOrUpdateValue(data.TTLSec, other.TTLSec, preserveKnown)
	data.Target = helper.KeepOrUpdateValue(data.Target, other.Target, preserveKnown)
	data.Priority = helper.KeepOrUpdateValue(data.Priority, other.Priority, preserveKnown)
	data.Weight = helper.KeepOrUpdateValue(data.Weight, other.Weight, preserveKnown)
	data.Port = helper.KeepOrUpdateValue(data.Port, other.Port, preserveKnown)
	data.Protocol = helper.KeepOrUpdateValue(data.Protocol, other.Protocol, preserveKnown)
	data.Service = helper.KeepOrUpdateValue(data.Service, other.Service, preserveKnown)
	data.Tag = helper.KeepOrUpdateValue(data.Tag, other.Tag, preserveKnown)
}
