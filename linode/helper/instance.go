package helper

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	RootPassMinimumCharacters     = 11
	RootPassMaximumCharacters     = 128
	DefaultFrameworkRebootTimeout = 600
)

var bootEvents = []linodego.EventAction{linodego.ActionLinodeBoot, linodego.ActionLinodeReboot}

// set bootConfig = 0 if using existing boot config
func RebootInstance(ctx context.Context, d *schema.ResourceData, linodeID int,
	meta interface{}, bootConfig int,
) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(GetDeadlineSeconds(ctx, d))*time.Second,
	)
	defer cancel()
	return diag.FromErr(rebootInstance(ctx, linodeID, &client, bootConfig))
}

func FrameworkRebootInstance(
	ctx context.Context,
	linodeID int,
	client *linodego.Client,
	bootConfig int,
) fwdiag.Diagnostics {
	var diags fwdiag.Diagnostics
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(
			ctx,
			time.Duration(DefaultFrameworkRebootTimeout)*time.Second,
		)
		defer cancel()
	}
	err := rebootInstance(ctx, linodeID, client, bootConfig)
	if err != nil {
		diags.AddError("Failed to Reboot Instance", err.Error())
	}
	return diags
}

func rebootInstance(
	ctx context.Context,
	entityID int,
	client *linodego.Client,
	bootConfig int,
) error {
	ctx = SetLogFieldBulk(ctx, map[string]any{
		"linode_id": entityID,
		"config_id": bootConfig,
	})
	instance, err := client.GetInstance(ctx, entityID)
	if err != nil {
		return fmt.Errorf("Error fetching data about the current linode: %s", err)
	}

	if instance.Status != linodego.InstanceRunning {
		tflog.Info(ctx, "Instance is not running")
		return nil
	}

	tflog.Info(ctx, "Rebooting instance")

	p, err := client.NewEventPoller(ctx, entityID, linodego.EntityLinode, linodego.ActionLinodeReboot)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	err = client.RebootInstance(ctx, instance.ID, bootConfig)
	if err != nil {
		return fmt.Errorf("Error rebooting Instance [%d]: %s", instance.ID, err)
	}

	deadlineSeconds := 600
	if deadline, ok := ctx.Deadline(); ok {
		deadlineSeconds = int(time.Until(deadline).Seconds())
	}
	_, err = p.WaitForFinished(ctx, deadlineSeconds)
	if err != nil {
		return fmt.Errorf("Error waiting for Instance [%d] to finish rebooting: %s", instance.ID, err)
	}

	if deadline, ok := ctx.Deadline(); ok {
		deadlineSeconds = int(time.Until(deadline).Seconds())
	}
	if _, err = client.WaitForInstanceStatus(
		ctx, instance.ID, linodego.InstanceRunning, deadlineSeconds,
	); err != nil {
		return fmt.Errorf("Timed-out waiting for Linode instance [%d] to boot: %s", instance.ID, err)
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

	// Special case for instances booted into rescue mode
	if events[0].SecondaryEntity == nil {
		return 0, nil
	}

	return int(events[0].SecondaryEntity.ID.(float64)), nil
}

func FrameworkCreateRandomRootPassword(diags *fwdiag.Diagnostics) string {
	rootPass, err := CreateRandomRootPassword()
	if err != nil {
		diags.AddError("Failed to Generate Random Root Password", err.Error())
		return ""
	}
	return rootPass
}

func CreateRandomRootPassword() (string, error) {
	rawRootPass := make([]byte, 50)
	_, err := rand.Read(rawRootPass)
	if err != nil {
		return "", fmt.Errorf("Failed to generate random password: %w", err)
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

	result := &linodego.VPCIPv4{
		VPC: vpcAddress,
	}

	if nat1To1 != "" {
		result.NAT1To1 = &nat1To1
	}

	return result
}

func ExpandConfigInterface(ifaceMap map[string]interface{}) linodego.InstanceConfigInterfaceCreateOptions {
	purpose := linodego.ConfigInterfacePurpose(ifaceMap["purpose"].(string))
	result := linodego.InstanceConfigInterfaceCreateOptions{
		Purpose: purpose,
		Primary: ifaceMap["primary"].(bool),
	}

	if purpose == linodego.InterfacePurposeVLAN {
		if ifaceMap["ipam_address"] != nil {
			result.IPAMAddress = ifaceMap["ipam_address"].(string)
		}

		if ifaceMap["label"] != nil {
			result.Label = ifaceMap["label"].(string)
		}
	}

	if purpose == linodego.InterfacePurposeVPC {
		if ifaceMap["subnet_id"] != nil {
			if subnetId := ifaceMap["subnet_id"].(int); subnetId != 0 {
				result.SubnetID = &subnetId
			}
		}

		if ifaceMap["ipv4"] != nil {
			if ipv4 := ifaceMap["ipv4"].([]any); len(ipv4) > 0 {
				result.IPv4 = ExpandInterfaceIPv4(ipv4[0])
			}
		}
		if ifaceMap["ip_ranges"] != nil {
			// this is for keep result.IPRanges as a nil value rather than a value of empty slice
			// when there is not a range.
			if ranges := ifaceMap["ip_ranges"].([]interface{}); len(ranges) > 0 {
				result.IPRanges = ExpandStringList(ranges)
			}
		}
	}

	return result
}

func ExpandConfigInterfaces(ctx context.Context, ifaces []any) []linodego.InstanceConfigInterfaceCreateOptions {
	result := make([]linodego.InstanceConfigInterfaceCreateOptions, len(ifaces))

	for i, iface := range ifaces {
		ifaceMap := iface.(map[string]any)
		result[i] = ExpandConfigInterface(ifaceMap)
	}

	return result
}

func FlattenInterfaceIPv4(ipv4 *linodego.VPCIPv4) []map[string]any {
	if ipv4 == nil {
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

// BootInstanceSync boots the instance with the given ID and waits for the operation to
// complete before returning.
func BootInstanceSync(
	ctx context.Context,
	client *linodego.Client,
	instanceID,
	configID,
	deadlineSeconds int,
) error {
	ctx = SetLogFieldBulk(
		ctx,
		map[string]any{
			"instance_id": instanceID,
			"config_id":   configID,
		},
	)

	tflog.Info(ctx, "Booting instance")

	p, err := client.NewEventPoller(ctx, instanceID, linodego.EntityLinode, linodego.ActionLinodeBoot)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	tflog.Debug(ctx, "client.BootInstance(...)")

	if err := client.BootInstance(ctx, instanceID, configID); err != nil {
		return fmt.Errorf("failed to boot instance: %s", err)
	}

	tflog.Debug(ctx, "Waiting for instance boot to finish")

	if _, err := p.WaitForFinished(ctx, deadlineSeconds); err != nil {
		return fmt.Errorf("failed to wait for instance boot: %s", err)
	}

	tflog.Debug(ctx, "Instance has finished booting")

	return nil
}

// ShutDownInstanceSync shuts down the instance with the given ID and waits for the operation to
// complete before returning.
func ShutDownInstanceSync(
	ctx context.Context,
	client *linodego.Client,
	instanceID,
	deadlineSeconds int,
) error {
	ctx = tflog.SetField(ctx, "instance_id", instanceID)

	tflog.Info(ctx, "Shutting down instance")

	p, err := client.NewEventPoller(ctx, instanceID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
	if err != nil {
		return fmt.Errorf("failed to initialize event poller: %s", err)
	}

	tflog.Debug(ctx, "client.ShutdownInstance(...)")

	if err := client.ShutdownInstance(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to shutdown instance: %s", err)
	}

	tflog.Debug(ctx, "Waiting for instance shutdown to finish")

	if _, err := p.WaitForFinished(ctx, deadlineSeconds); err != nil {
		return fmt.Errorf("failed to wait for instance shutdown: %s", err)
	}

	tflog.Debug(ctx, "Instance has finished shutting down")

	return nil
}

// WaitForInstanceNonTransientStatus waits for the instance with the given ID to enter
// a non-transient status (e.g. running, offline), and returns the final status of the instance.
func WaitForInstanceNonTransientStatus(
	ctx context.Context,
	client *linodego.Client,
	linodeID int,
	timeoutSeconds int,
) (linodego.InstanceStatus, error) {
	instance, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return "", fmt.Errorf("failed to get instance: %w", err)
	}

	var targetStatus linodego.InstanceStatus

	switch instance.Status {
	case linodego.InstanceBooting, linodego.InstanceRebooting:
		targetStatus = linodego.InstanceRunning

	case linodego.InstanceShuttingDown:
		targetStatus = linodego.InstanceOffline

	case linodego.InstanceRunning, linodego.InstanceOffline:
		// Instance is offline, nothing to do here
		return instance.Status, nil

	default:
		return "", fmt.Errorf("cannot wait for instance to exit transient status %s", instance.Status)
	}

	instance, err = client.WaitForInstanceStatus(
		ctx,
		instance.ID,
		targetStatus,
		timeoutSeconds,
	)
	if err != nil {
		return "", fmt.Errorf("failed to wait for instance to reach status %s: %w", targetStatus, err)
	}

	return instance.Status, nil
}
