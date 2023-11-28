package helper

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	RootPassMinimumCharacters = 11
	RootPassMaximumCharacters = 128
)

var bootEvents = []linodego.EventAction{linodego.ActionLinodeBoot, linodego.ActionLinodeReboot}

// set bootConfig = 0 if using existing boot config
func RebootInstance(ctx context.Context, d *schema.ResourceData, entityID int,
	meta interface{}, bootConfig int,
) diag.Diagnostics {
	ctx = SetLogFieldBulk(ctx, map[string]any{
		"linode_id": entityID,
		"config_id": bootConfig,
	})

	client := meta.(*ProviderMeta).Client
	instance, err := client.GetInstance(ctx, entityID)
	if err != nil {
		return diag.Errorf("Error fetching data about the current linode: %s", err)
	}

	if instance.Status != linodego.InstanceRunning {
		tflog.Info(ctx, "Instance is not running")
		return nil
	}

	tflog.Info(ctx, "Rebooting instance")

	p, err := client.NewEventPoller(ctx, entityID, linodego.EntityLinode, linodego.ActionLinodeReboot)
	if err != nil {
		return diag.Errorf("failed to initialize event poller: %s", err)
	}

	err = client.RebootInstance(ctx, instance.ID, bootConfig)

	if err != nil {
		return diag.Errorf("Error rebooting Instance [%d]: %s", instance.ID, err)
	}
	_, err = p.WaitForFinished(ctx, GetDeadlineSeconds(ctx, d))
	if err != nil {
		return diag.Errorf("Error waiting for Instance [%d] to finish rebooting: %s", instance.ID, err)
	}

	if _, err = client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceRunning, GetDeadlineSeconds(ctx, d),
	); err != nil {
		return diag.Errorf("Timed-out waiting for Linode instance [%d] to boot: %s", instance.ID, err)
	}

	tflog.Debug(ctx, "Instance has finished rebooting")

	return nil
}

// GetDeadlineSeconds gets the seconds remaining until deadline is met.
func GetDeadlineSeconds(ctx context.Context, d *schema.ResourceData) int {
	duration := d.Timeout(schema.TimeoutUpdate)
	if deadline, ok := ctx.Deadline(); ok {
		duration = time.Until(deadline)
	}
	return int(duration.Seconds())
}

// IsInstanceInBootedState checks whether an instance is in a booted or booting state
func IsInstanceInBootedState(status linodego.InstanceStatus) bool {
	// For diffing purposes, transition states need to be treated as
	// booted == true. This is because these statuses will eventually
	// result in a powered on Linode.
	return status == linodego.InstanceRunning ||
		status == linodego.InstanceRebooting ||
		status == linodego.InstanceBooting
}

// GetCurrentBootedConfig gets the config a linode instance is current booted to
func GetCurrentBootedConfig(ctx context.Context, client *linodego.Client, instID int) (int, error) {
	inst, err := client.GetInstance(ctx, instID)
	if err != nil {
		return 0, err
	}

	// Valid exit condition where no config is booted
	if !IsInstanceInBootedState(inst.Status) {
		return 0, nil
	}

	filter := map[string]any{
		"entity.id":   instID,
		"entity.type": linodego.EntityLinode,
		"+or":         []map[string]any{},
		"+order_by":   "created",
		"+order":      "desc",
	}

	for _, v := range bootEvents {
		filter["+or"] = append(filter["+or"].([]map[string]any), map[string]any{"action": v})
	}

	filterBytes, err := json.Marshal(filter)
	if err != nil {
		return 0, err
	}

	events, err := client.ListEvents(ctx, &linodego.ListOptions{
		Filter: string(filterBytes),
	})
	if err != nil {
		return 0, err
	}

	if len(events) < 1 {
		// This is a valid exit case
		return 0, nil
	}

	return int(events[0].SecondaryEntity.ID.(float64)), nil
}

func CreateRandomRootPassword() (string, error) {
	rawRootPass := make([]byte, 50)
	_, err := rand.Read(rawRootPass)
	if err != nil {
		return "", fmt.Errorf("Failed to generate random password")
	}
	rootPass := base64.StdEncoding.EncodeToString(rawRootPass)
	return rootPass, nil
}

func ExpandInterfaceIPv4(ipv4 any) *linodego.VPCIPv4 {
	IPv4 := ipv4.(map[string]any)
	vpcAddress := IPv4["vpc"].(string)
	nat1To1 := IPv4["nat_1_1"].(string)
	if vpcAddress == "" && nat1To1 == "" {
		return nil
	}
	return &linodego.VPCIPv4{
		VPC:     vpcAddress,
		NAT1To1: nat1To1,
	}
}

func ExpandConfigInterface(ifaceMap map[string]interface{}) linodego.InstanceConfigInterfaceCreateOptions {
	result := linodego.InstanceConfigInterfaceCreateOptions{
		Purpose:     linodego.ConfigInterfacePurpose(ifaceMap["purpose"].(string)),
		Label:       ifaceMap["label"].(string),
		IPAMAddress: ifaceMap["ipam_address"].(string),
		Primary:     ifaceMap["primary"].(bool),
	}
	if ifaceMap["subnet_id"] != nil {
		subnet_id := ifaceMap["subnet_id"].(int)
		if subnet_id != 0 {
			result.SubnetID = &subnet_id
		}
	}

	if ifaceMap["ipv4"] != nil {
		ipv4 := ifaceMap["ipv4"].([]any)
		if len(ipv4) > 0 {
			result.IPv4 = ExpandInterfaceIPv4(ipv4[0])
		}
	}
	if ifaceMap["ip_ranges"] != nil {
		result.IPRanges = ExpandStringList(ifaceMap["ip_ranges"].([]interface{}))
	}

	return result
}

func ExpandInterfaces(ctx context.Context, ifaces []any) []linodego.InstanceConfigInterfaceCreateOptions {
	result := make([]linodego.InstanceConfigInterfaceCreateOptions, len(ifaces))

	for i, iface := range ifaces {
		ifaceMap := iface.(map[string]any)
		result[i] = ExpandConfigInterface(ifaceMap)
	}

	return result
}

func FlattenInterfaceIPv4(ipv4 linodego.VPCIPv4) []map[string]any {
	if ipv4.NAT1To1 == "" && ipv4.VPC == "" {
		return nil
	}
	return []map[string]any{
		{
			"vpc":     ipv4.VPC,
			"nat_1_1": ipv4.NAT1To1,
		},
	}
}

func FlattenInterface(iface linodego.InstanceConfigInterface) map[string]any {
	return map[string]any{
		"purpose":      iface.Purpose,
		"ipam_address": iface.IPAMAddress,
		"label":        iface.Label,
		"id":           iface.ID,
		"vpc_id":       iface.VPCID,
		"subnet_id":    iface.SubnetID,
		"primary":      iface.Primary,
		"active":       iface.Active,
		"ip_ranges":    iface.IPRanges,
		"ipv4":         FlattenInterfaceIPv4(iface.IPv4),
	}
}

func FlattenInterfaces(interfaces []linodego.InstanceConfigInterface) []map[string]any {
	result := make([]map[string]any, len(interfaces))
	for i, iface := range interfaces {
		result[i] = FlattenInterface(iface)
	}
	return result
}
