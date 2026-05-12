package instancenetworking

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DataSourceModel struct {
	LinodeID types.Int64  `tfsdk:"linode_id"`
	IPV4     types.List   `tfsdk:"ipv4"`
	IPV6     types.List   `tfsdk:"ipv6"`
	ID       types.String `tfsdk:"id"`
}

func (data *DataSourceModel) parseInstanceIPAddressResponse(
	ip *linodego.InstanceIPAddressResponse, diags *diag.Diagnostics,
) {
	ipv4 := flattenIPv4(ip.IPv4, diags)
	if diags.HasError() {
		return
	}

	data.IPV4 = *ipv4

	ipv6 := flattenIPv6(ip.IPv6, diags)
	if diags.HasError() {
		return
	}

	data.IPV6 = *ipv6

	id, err := json.Marshal(ip)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return
	}

	data.ID = types.StringValue(string(id))
}

func flattenIPv4(network *linodego.InstanceIPv4Response, diags *diag.Diagnostics) *basetypes.ListValue {
	result := make(map[string]attr.Value)

	result["private"] = helper.GenericSliceToList(network.Private, networkObjectType, flattenIP, diags)
	if diags.HasError() {
		return nil
	}

	result["public"] = helper.GenericSliceToList(network.Public, networkObjectType, flattenIP, diags)
	if diags.HasError() {
		return nil
	}

	result["reserved"] = helper.GenericSliceToList(network.Reserved, networkObjectType, flattenIP, diags)

	if diags.HasError() {
		return nil
	}

	result["shared"] = helper.GenericSliceToList(network.Shared, networkObjectType, flattenIP, diags)
	if diags.HasError() {
		return nil
	}

	result["vpc"] = helper.GenericSliceToList(network.VPC, vpcNetworkObjectType, flattenVPCIP, diags)
	if diags.HasError() {
		return nil
	}

	resultList := helper.MapToSingleObjList(ipv4ObjectType, result, diags)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func flattenIPv6(network *linodego.InstanceIPv6Response, diags *diag.Diagnostics) *basetypes.ListValue {
	result := make(map[string]attr.Value)

	global := helper.GenericSliceToList(network.Global, globalObjectType, flattenIPV6Range, diags)

	link_local, newDiags := flattenIP(network.LinkLocal)
	if newDiags.HasError() {
		diags.Append(newDiags...)
		return nil
	}

	slaac, newDiags := flattenIP(network.SLAAC)
	if newDiags.HasError() {
		diags.Append(newDiags...)
		return nil
	}

	result["global"] = global
	result["link_local"] = link_local
	result["slaac"] = slaac
	result["vpc"] = helper.GenericSliceToList(network.VPC, vpcNetworkObjectType, flattenVPCIPByValue, diags)
	if diags.HasError() {
		return nil
	}

	obj, newDiags := types.ObjectValue(ipv6ObjectType.AttrTypes, result)
	if newDiags.HasError() {
		diags.Append(newDiags...)
		return nil
	}

	resultList := helper.GenericSliceToList(
		[]attr.Value{obj}, ipv6ObjectType, helper.FwValueEchoConverter(), diags,
	)
	if diags.HasError() {
		return nil
	}

	return &resultList
}

func FlattenIPVPCNAT1To1(data *linodego.InstanceIPNAT1To1) (basetypes.ObjectValue, diag.Diagnostics) {
	if data == nil {
		return types.ObjectNull(VPCNAT1To1Type.AttrTypes), nil
	}

	result := map[string]attr.Value{
		"address":   types.StringValue(data.Address),
		"vpc_id":    types.Int64Value(int64(data.VPCID)),
		"subnet_id": types.Int64Value(int64(data.SubnetID)),
	}

	obj, diag := types.ObjectValue(VPCNAT1To1Type.AttrTypes, result)
	if diag.HasError() {
		return types.ObjectNull(VPCNAT1To1Type.AttrTypes), diag
	}

	return obj, nil
}

func flattenVPCIP(vpc *linodego.VPCIP) (*basetypes.ObjectValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)
	var diags diag.Diagnostics

	result["address"] = types.StringPointerValue(vpc.Address)
	result["address_range"] = types.StringPointerValue(vpc.AddressRange)
	result["active"] = types.BoolValue(vpc.Active)
	result["vpc_id"] = types.Int64Value(int64(vpc.VPCID))
	result["subnet_id"] = types.Int64Value(int64(vpc.SubnetID))
	result["config_id"] = types.Int64Value(int64(vpc.ConfigID))
	result["interface_id"] = types.Int64Value(int64(vpc.ConfigID))

	result["gateway"] = types.StringValue(vpc.Gateway)
	result["prefix"] = types.Int64Value(int64(vpc.Prefix))
	result["region"] = types.StringValue(vpc.Region)
	result["subnet_mask"] = types.StringValue(vpc.SubnetMask)
	result["linode_id"] = types.Int64Value(int64(vpc.LinodeID))
	result["nat_1_1"] = types.StringPointerValue(vpc.NAT1To1)

	result["ipv6_range"] = types.StringPointerValue(vpc.IPv6Range)
	result["ipv6_is_public"] = types.BoolPointerValue(vpc.IPv6IsPublic)
	result["ipv6_addresses"] = helper.GenericSliceToList(
		vpc.IPv6Addresses,
		vpcIPv6AddressObjectType,
		flattenVPCIPIPv6Address,
		&diags,
	)
	if diags.HasError() {
		return nil, diags
	}

	obj, newDiags := types.ObjectValue(vpcNetworkObjectType.AttrTypes, result)
	if newDiags.HasError() {
		return nil, newDiags
	}

	return &obj, nil
}

func flattenIP(network *linodego.InstanceIP) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	vpcNAT1To1, d := FlattenIPVPCNAT1To1(network.VPCNAT1To1)
	if d.HasError() {
		return nil, d
	}

	result["vpc_nat_1_1"] = vpcNAT1To1

	result["address"] = types.StringValue(network.Address)
	result["gateway"] = types.StringValue(network.Gateway)
	result["prefix"] = types.Int64Value(int64(network.Prefix))
	result["rdns"] = types.StringValue(network.RDNS)
	result["region"] = types.StringValue(network.Region)
	result["subnet_mask"] = types.StringValue(network.SubnetMask)
	result["type"] = types.StringValue(string(network.Type))
	result["public"] = types.BoolValue(network.Public)
	result["linode_id"] = types.Int64Value(int64(network.LinodeID))

	if network.InterfaceID != nil {
		result["interface_id"] = types.Int64Value(int64(*network.InterfaceID))
	} else {
		result["interface_id"] = types.Int64Null()
	}

	result["reserved"] = types.BoolValue(network.Reserved)

	tags := helper.StringSliceToFrameworkValueSlice(network.Tags)
	tagsSet, tagsDiags := types.SetValue(types.StringType, tags)
	if tagsDiags.HasError() {
		return nil, tagsDiags
	}
	result["tags"] = tagsSet

	assignedEntity, assignedEntityDiags := FlattenAssignedEntity(network.AssignedEntity)
	if assignedEntityDiags.HasError() {
		return nil, assignedEntityDiags
	}
	result["assigned_entity"] = assignedEntity

	obj, d := types.ObjectValue(networkObjectType.AttrTypes, result)
	if d.HasError() {
		return nil, d
	}

	return &obj, nil
}

func FlattenAssignedEntity(entity *linodego.ReservedIPAssignedEntity) (basetypes.ObjectValue, diag.Diagnostics) {
	if entity == nil {
		return types.ObjectNull(AssignedEntityObjectType.AttrTypes), nil
	}

	result := map[string]attr.Value{
		"id":    types.Int64Value(int64(entity.ID)),
		"label": types.StringValue(entity.Label),
		"type":  types.StringValue(entity.Type),
		"url":   types.StringValue(entity.URL),
	}

	return types.ObjectValue(AssignedEntityObjectType.AttrTypes, result)
}

func flattenVPCIPByValue(vpc linodego.VPCIP) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	return flattenVPCIP(&vpc)
}

func flattenVPCIPIPv6Address(addr linodego.VPCIPIPv6Address) (
	basetypes.ObjectValue, diag.Diagnostics,
) {
	return types.ObjectValue(vpcIPv6AddressObjectType.AttrTypes, map[string]attr.Value{
		"slaac_address": types.StringValue(addr.SLAACAddress),
	})
}

func flattenIPV6Range(network linodego.IPv6Range) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["prefix"] = types.Int64Value(int64(network.Prefix))
	result["range"] = types.StringValue(network.Range)
	result["region"] = types.StringValue(network.Region)
	result["route_target"] = types.StringValue(network.RouteTarget)

	obj, diag := types.ObjectValue(globalObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	return &obj, nil
}
