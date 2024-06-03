package nbconfig

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	ProxyProtocol  types.String `tfsdk:"proxy_protocol"`
	Port           types.Int64  `tfsdk:"port"`
	CheckInterval  types.Int64  `tfsdk:"check_interval"`
	CheckTimeout   types.Int64  `tfsdk:"check_timeout"`
	CheckAttempts  types.Int64  `tfsdk:"check_attempts"`
	Algorithm      types.String `tfsdk:"algorithm"`
	Stickiness     types.String `tfsdk:"stickiness"`
	Check          types.String `tfsdk:"check"`
	CheckPath      types.String `tfsdk:"check_path"`
	CheckBody      types.String `tfsdk:"check_body"`
	CheckPassive   types.Bool   `tfsdk:"check_passive"`
	CipherSuite    types.String `tfsdk:"cipher_suite"`
	SSLCommonName  types.String `tfsdk:"ssl_commonname"`
	SSLFingerprint types.String `tfsdk:"ssl_fingerprint"`
	NodesStatus    types.List   `tfsdk:"node_status"`
}

func (data *DataSourceModel) ParseNodebalancerConfig(
	config *linodego.NodeBalancerConfig,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(config.ID))
	data.NodeBalancerID = types.Int64Value(int64(config.NodeBalancerID))
	data.Algorithm = types.StringValue(string(config.Algorithm))
	data.Stickiness = types.StringValue(string(config.Stickiness))
	data.Check = types.StringValue(string(config.Check))
	data.CheckAttempts = types.Int64Value(int64(config.CheckAttempts))
	data.CheckBody = types.StringValue(config.CheckBody)
	data.CheckInterval = types.Int64Value(int64(config.CheckInterval))
	data.CheckTimeout = types.Int64Value(int64(config.CheckTimeout))
	data.CheckPassive = types.BoolValue(config.CheckPassive)
	data.CheckPath = types.StringValue(config.CheckPath)
	data.CipherSuite = types.StringValue(string(config.CipherSuite))
	data.Port = types.Int64Value(int64(config.Port))
	data.Protocol = types.StringValue(string(config.Protocol))
	data.ProxyProtocol = types.StringValue(string(config.ProxyProtocol))
	data.SSLFingerprint = types.StringValue(config.SSLFingerprint)
	data.SSLCommonName = types.StringValue(config.SSLCommonName)

	nodeStatus, diags := flattenNodeStatus(config.NodesStatus)
	if diags.HasError() {
		return diags
	}
	data.NodesStatus = *nodeStatus
	return nil
}

func flattenNodeStatus(
	nodesStatus *linodego.NodeBalancerNodeStatus,
) (*types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make(map[string]attr.Value)

	result["up"] = types.Int64Value(int64(nodesStatus.Up))
	result["down"] = types.Int64Value(int64(nodesStatus.Down))

	resultList := helper.MapToSingleObjList(statusObjectType, result, &diags)
	return &resultList, diags
}

type ResourceModelV0 struct {
	ID             types.String `tfsdk:"id"`
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	ProxyProtocol  types.String `tfsdk:"proxy_protocol"`
	Port           types.Int64  `tfsdk:"port"`
	CheckInterval  types.Int64  `tfsdk:"check_interval"`
	CheckTimeout   types.Int64  `tfsdk:"check_timeout"`
	CheckAttempts  types.Int64  `tfsdk:"check_attempts"`
	Algorithm      types.String `tfsdk:"algorithm"`
	Stickiness     types.String `tfsdk:"stickiness"`
	Check          types.String `tfsdk:"check"`
	CheckPath      types.String `tfsdk:"check_path"`
	CheckBody      types.String `tfsdk:"check_body"`
	CheckPassive   types.Bool   `tfsdk:"check_passive"`
	CipherSuite    types.String `tfsdk:"cipher_suite"`
	SSLCommonName  types.String `tfsdk:"ssl_commonname"`
	SSLFingerprint types.String `tfsdk:"ssl_fingerprint"`
	NodesStatus    types.Map    `tfsdk:"node_status"`
	SSLCert        types.String `tfsdk:"ssl_cert"`
	SSLKey         types.String `tfsdk:"ssl_key"`
}

type ResourceModelV1 struct {
	ID             types.String `tfsdk:"id"`
	NodeBalancerID types.Int64  `tfsdk:"nodebalancer_id"`
	Protocol       types.String `tfsdk:"protocol"`
	ProxyProtocol  types.String `tfsdk:"proxy_protocol"`
	Port           types.Int64  `tfsdk:"port"`
	CheckInterval  types.Int64  `tfsdk:"check_interval"`
	CheckTimeout   types.Int64  `tfsdk:"check_timeout"`
	CheckAttempts  types.Int64  `tfsdk:"check_attempts"`
	Algorithm      types.String `tfsdk:"algorithm"`
	Stickiness     types.String `tfsdk:"stickiness"`
	Check          types.String `tfsdk:"check"`
	CheckPath      types.String `tfsdk:"check_path"`
	CheckBody      types.String `tfsdk:"check_body"`
	CheckPassive   types.Bool   `tfsdk:"check_passive"`
	CipherSuite    types.String `tfsdk:"cipher_suite"`
	SSLCommonName  types.String `tfsdk:"ssl_commonname"`
	SSLFingerprint types.String `tfsdk:"ssl_fingerprint"`
	NodesStatus    types.List   `tfsdk:"node_status"`
	SSLCert        types.String `tfsdk:"ssl_cert"`
	SSLKey         types.String `tfsdk:"ssl_key"`
}

func (data *ResourceModelV1) FlattenNodeBalancerConfig(
	config *linodego.NodeBalancerConfig, preserveKnown bool,
) diag.Diagnostics {
	data.ID = helper.KeepOrUpdateString(data.ID, strconv.Itoa(config.ID), preserveKnown)
	data.NodeBalancerID = helper.KeepOrUpdateInt64(data.NodeBalancerID, int64(config.NodeBalancerID), preserveKnown)
	data.Algorithm = helper.KeepOrUpdateString(data.Algorithm, string(config.Algorithm), preserveKnown)
	data.Stickiness = helper.KeepOrUpdateString(data.Stickiness, string(config.Stickiness), preserveKnown)
	data.Check = helper.KeepOrUpdateString(data.Check, string(config.Check), preserveKnown)
	data.CheckAttempts = helper.KeepOrUpdateInt64(data.CheckAttempts, int64(config.CheckAttempts), preserveKnown)
	data.CheckBody = helper.KeepOrUpdateString(data.CheckBody, config.CheckBody, preserveKnown)
	data.CheckInterval = helper.KeepOrUpdateInt64(data.CheckInterval, int64(config.CheckInterval), preserveKnown)
	data.CheckTimeout = helper.KeepOrUpdateInt64(data.CheckTimeout, int64(config.CheckTimeout), preserveKnown)
	data.CheckPassive = helper.KeepOrUpdateBool(data.CheckPassive, config.CheckPassive, preserveKnown)
	data.CheckPath = helper.KeepOrUpdateString(data.CheckPath, config.CheckPath, preserveKnown)
	data.CipherSuite = helper.KeepOrUpdateString(data.CipherSuite, string(config.CipherSuite), preserveKnown)
	data.Port = helper.KeepOrUpdateInt64(data.Port, int64(config.Port), preserveKnown)
	data.Protocol = helper.KeepOrUpdateString(data.Protocol, string(config.Protocol), preserveKnown)
	data.ProxyProtocol = helper.KeepOrUpdateString(data.ProxyProtocol, string(config.ProxyProtocol), preserveKnown)
	data.SSLFingerprint = helper.KeepOrUpdateString(data.SSLFingerprint, config.SSLFingerprint, preserveKnown)
	data.SSLCommonName = helper.KeepOrUpdateString(data.SSLCommonName, config.SSLCommonName, preserveKnown)
	// SSLCert and SSLKey are not included because they are
	// neither computed nor returned from the GET API call.

	nodeStatus, diags := flattenNodeStatus(config.NodesStatus)
	if diags.HasError() {
		return diags
	}
	data.NodesStatus = helper.KeepOrUpdateValue(data.NodesStatus, *nodeStatus, preserveKnown)

	return nil
}

func (data *ResourceModelV1) CopyFrom(other ResourceModelV1, preserveKnown bool) {
	data.ID = helper.KeepOrUpdateValue(data.ID, other.ID, preserveKnown)
	data.NodeBalancerID = helper.KeepOrUpdateValue(data.NodeBalancerID, other.NodeBalancerID, preserveKnown)
	data.Algorithm = helper.KeepOrUpdateValue(data.Algorithm, other.Algorithm, preserveKnown)
	data.Stickiness = helper.KeepOrUpdateValue(data.Stickiness, other.Stickiness, preserveKnown)
	data.Check = helper.KeepOrUpdateValue(data.Check, other.Check, preserveKnown)
	data.CheckAttempts = helper.KeepOrUpdateValue(data.CheckAttempts, other.CheckAttempts, preserveKnown)
	data.CheckBody = helper.KeepOrUpdateValue(data.CheckBody, other.CheckBody, preserveKnown)
	data.CheckInterval = helper.KeepOrUpdateValue(data.CheckInterval, other.CheckInterval, preserveKnown)
	data.CheckTimeout = helper.KeepOrUpdateValue(data.CheckTimeout, other.CheckTimeout, preserveKnown)
	data.CheckPassive = helper.KeepOrUpdateValue(data.CheckPassive, other.CheckPassive, preserveKnown)
	data.CheckPath = helper.KeepOrUpdateValue(data.CheckPath, other.CheckPath, preserveKnown)
	data.CipherSuite = helper.KeepOrUpdateValue(data.CipherSuite, other.CipherSuite, preserveKnown)
	data.Port = helper.KeepOrUpdateValue(data.Port, other.Port, preserveKnown)
	data.Protocol = helper.KeepOrUpdateValue(data.Protocol, other.Protocol, preserveKnown)
	data.ProxyProtocol = helper.KeepOrUpdateValue(data.ProxyProtocol, other.ProxyProtocol, preserveKnown)
	data.SSLFingerprint = helper.KeepOrUpdateValue(data.SSLFingerprint, other.SSLFingerprint, preserveKnown)
	data.SSLCommonName = helper.KeepOrUpdateValue(data.SSLCommonName, other.SSLCommonName, preserveKnown)
	data.SSLCert = helper.KeepOrUpdateValue(data.SSLCert, other.SSLCert, preserveKnown)
	data.SSLKey = helper.KeepOrUpdateValue(data.SSLKey, other.SSLKey, preserveKnown)
	data.NodesStatus = helper.KeepOrUpdateValue(data.NodesStatus, other.NodesStatus, preserveKnown)
}

func (v1 *ResourceModelV1) UpgradeFromV0(
	ctx context.Context, v0 ResourceModelV0,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if v0.NodesStatus.IsNull() || v0.NodesStatus.IsUnknown() {
		return diags
	}

	v1.ID = v0.ID
	v1.NodeBalancerID = v0.NodeBalancerID
	v1.Protocol = v0.Protocol
	v1.ProxyProtocol = v0.ProxyProtocol
	v1.Port = v0.Port
	v1.CheckInterval = v0.CheckInterval
	v1.CheckTimeout = v0.CheckTimeout
	v1.CheckAttempts = v0.CheckAttempts
	v1.Algorithm = v0.Algorithm
	v1.Stickiness = v0.Stickiness
	v1.Check = v0.Check
	v1.CheckPath = v0.CheckPath
	v1.CheckBody = v0.CheckBody
	v1.CheckPassive = v0.CheckPassive
	v1.CipherSuite = v0.CipherSuite
	v1.SSLCommonName = v0.SSLCommonName
	v1.SSLFingerprint = v0.SSLFingerprint
	v1.SSLCert = v0.SSLCert
	v1.SSLKey = v0.SSLKey

	nodesStatusV0 := make(map[string]types.String, len(v0.NodesStatus.Elements()))
	newDiags := v0.NodesStatus.ElementsAs(ctx, &nodesStatusV0, false)
	diags.Append(newDiags...)
	if diags.HasError() {
		return diags
	}

	nodesStatusV1 := make(map[string]attr.Value, len(nodesStatusV0))

	oldDown := nodesStatusV0["down"].ValueString()
	oldUp := nodesStatusV0["up"].ValueString()

	if oldDown != "" {
		down, err := strconv.Atoi(oldDown)
		if err != nil {
			diags.AddError("Failed to Convert 'down' to Int64", err.Error())
			return diags
		}
		nodesStatusV1["down"] = types.Int64Value(int64(down))
	} else {
		nodesStatusV1["down"] = types.Int64Value(0)
	}

	if oldUp != "" {
		up, err := strconv.Atoi(oldUp)
		if err != nil {
			diags.AddError("Failed to Convert 'up' to Int64", err.Error())
			return diags
		}
		nodesStatusV1["up"] = types.Int64Value(int64(up))
	} else {
		nodesStatusV1["up"] = types.Int64Value(0)
	}

	v1.NodesStatus = helper.MapToSingleObjList(
		NodeStatusTypeV1, nodesStatusV1, &diags,
	)

	return diags
}

func (data *ResourceModelV1) GetNodeBalancerConfigCreateOptions(
	ctx context.Context, diags *diag.Diagnostics,
) *linodego.NodeBalancerConfigCreateOptions {
	checkAttempts := helper.FrameworkSafeInt64ToInt(data.CheckAttempts.ValueInt64(), diags)
	checkInterval := helper.FrameworkSafeInt64ToInt(data.CheckInterval.ValueInt64(), diags)
	checkTimeout := helper.FrameworkSafeInt64ToInt(data.CheckTimeout.ValueInt64(), diags)
	port := helper.FrameworkSafeInt64ToInt(data.Port.ValueInt64(), diags)
	if diags.HasError() {
		return nil
	}
	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Algorithm:     linodego.ConfigAlgorithm(data.Algorithm.ValueString()),
		Check:         linodego.ConfigCheck(data.Check.ValueString()),
		Stickiness:    linodego.ConfigStickiness(data.Stickiness.ValueString()),
		CheckAttempts: checkAttempts,
		CheckBody:     data.CheckBody.ValueString(),
		CheckInterval: checkInterval,
		CheckPath:     data.CheckPath.ValueString(),
		CheckTimeout:  checkTimeout,
		Port:          port,
		Protocol:      linodego.ConfigProtocol(strings.ToLower(data.Protocol.ValueString())),
		ProxyProtocol: linodego.ConfigProxyProtocol(data.ProxyProtocol.ValueString()),
		SSLCert:       data.SSLCert.ValueString(),
		SSLKey:        data.SSLKey.ValueString(),
	}

	if !data.CheckPassive.IsUnknown() {
		createOpts.CheckPassive = data.CheckPassive.ValueBoolPointer()
	}

	return &createOpts
}

func (data *ResourceModelV1) GetNodeBalancerConfigUpdateOptions(
	ctx context.Context, diags *diag.Diagnostics,
) *linodego.NodeBalancerConfigUpdateOptions {
	checkAttempts := helper.FrameworkSafeInt64ToInt(data.CheckAttempts.ValueInt64(), diags)
	checkInterval := helper.FrameworkSafeInt64ToInt(data.CheckInterval.ValueInt64(), diags)
	checkTimeout := helper.FrameworkSafeInt64ToInt(data.CheckTimeout.ValueInt64(), diags)
	port := helper.FrameworkSafeInt64ToInt(data.Port.ValueInt64(), diags)
	if diags.HasError() {
		return nil
	}

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		Algorithm:     linodego.ConfigAlgorithm(data.Algorithm.ValueString()),
		Check:         linodego.ConfigCheck(data.Check.ValueString()),
		Stickiness:    linodego.ConfigStickiness(data.Stickiness.ValueString()),
		CheckAttempts: checkAttempts,
		CheckBody:     data.CheckBody.ValueString(),
		CheckInterval: checkInterval,
		CheckPath:     data.CheckPath.ValueString(),
		CheckTimeout:  checkTimeout,
		Port:          port,
		Protocol:      linodego.ConfigProtocol(strings.ToLower(data.Protocol.ValueString())),
		ProxyProtocol: linodego.ConfigProxyProtocol(data.ProxyProtocol.ValueString()),
		SSLCert:       data.SSLCert.ValueString(),
		SSLKey:        data.SSLKey.ValueString(),
	}

	if !data.CheckPassive.IsUnknown() {
		updateOpts.CheckPassive = data.CheckPassive.ValueBoolPointer()
	}

	return &updateOpts
}
