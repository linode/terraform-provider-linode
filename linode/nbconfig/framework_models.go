package nbconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	NodebalancerId types.Int64  `tfsdk:"nodebalancer_id"`
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
	SSLCommonname  types.String `tfsdk:"ssl_commonname"`
	SSLFingerprint types.String `tfsdk:"ssl_fingerprint"`
	NodeStatus     types.List   `tfsdk:"node_status"`
}

func (data *DataSourceModel) ParseNodebalancerConfig(
	config *linodego.NodeBalancerConfig,
) diag.Diagnostics {
	data.ID = types.Int64Value(int64(config.ID))
	data.NodebalancerId = types.Int64Value(int64(config.NodeBalancerID))
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
	data.SSLCommonname = types.StringValue(config.SSLCommonName)

	nodeStatus, diags := parseNodeStatus(config.NodesStatus)
	if diags.HasError() {
		return diags
	}
	data.NodeStatus = *nodeStatus
	return nil
}

func parseNodeStatus(
	nodesStatus *linodego.NodeBalancerNodeStatus,
) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["up"] = types.Int64Value(int64(nodesStatus.Up))
	result["down"] = types.Int64Value(int64(nodesStatus.Down))

	nodeStatusObj, diags := types.ObjectValue(statusObjectType.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}

	resultList, diags := basetypes.NewListValue(
		statusObjectType,
		[]attr.Value{nodeStatusObj},
	)
	if diags.HasError() {
		return nil, diags
	}
	return &resultList, nil
}
