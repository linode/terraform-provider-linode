package instance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	linodediffs "github.com/linode/terraform-provider-linode/v2/linode/helper/customdiffs"
)

const (
	LinodeInstanceCreateTimeout = 15 * time.Minute
	LinodeInstanceUpdateTimeout = time.Hour
	LinodeInstanceDeleteTimeout = 10 * time.Minute
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		CustomizeDiff: customdiff.All(
			linodediffs.ComputedWithDefault("tags", []string{}),
			linodediffs.CaseInsensitiveSet("tags"),
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LinodeInstanceCreateTimeout),
			Update: schema.DefaultTimeout(LinodeInstanceUpdateTimeout),
			Delete: schema.DefaultTimeout(LinodeInstanceDeleteTimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_instance")

	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode instance ID %s as int: %s", d.Id(), err)
	}

	instance, err := client.GetInstance(ctx, id)
	if linodego.IsNotFound(err) {
		tflog.Warn(ctx, "Removing Linode Instance ID %q from state because it no longer exists")
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("failed to get instance: %s", err)
	}

	instanceNetwork, err := client.GetInstanceIPAddresses(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get instance networking: %s", err)
	}

	instanceDisks, err := client.ListInstanceDisks(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get instance disks: %s", err)
	}

	instanceConfigs, err := client.ListInstanceConfigs(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get instance configs: %s", err)
	}

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}
	d.Set("ipv4", ips)
	d.Set("ipv6", instance.IPv6)
	d.Set("shared_ipv4", instanceIPSliceToString(instanceNetwork.IPv4.Shared))

	public, private := instanceNetwork.IPv4.Public, instanceNetwork.IPv4.Private

	if len(public) > 0 {
		d.Set("ip_address", public[0].Address)

		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": public[0].Address,
		})
	}

	if len(private) > 0 {
		d.Set("private_ip", true)
		d.Set("private_ip_address", private[0].Address)
	} else {
		d.Set("private_ip", false)
	}

	d.Set("label", instance.Label)
	d.Set("status", instance.Status)
	d.Set("type", instance.Type)
	d.Set("region", instance.Region)
	d.Set("watchdog_enabled", instance.WatchdogEnabled)
	d.Set("group", instance.Group)
	d.Set("tags", instance.Tags)
	d.Set("booted", isInstanceBooted(instance))
	d.Set("host_uuid", instance.HostUUID)
	d.Set("has_user_data", instance.HasUserData)
	d.Set("lke_cluster_id", instance.LKEClusterID)
	d.Set("disk_encryption", instance.DiskEncryption)

	flatSpecs := flattenInstanceSpecs(*instance)
	flatAlerts := flattenInstanceAlerts(*instance)
	flatBackups := flattenInstanceBackups(*instance)

	d.Set("backups", flatBackups)
	d.Set("backups_enabled", instance.Backups.Enabled)

	d.Set("specs", flatSpecs)
	d.Set("alerts", flatAlerts)

	var placementGroupMap map[string]interface{}
	flattenedGroups := flattenInstancePlacementGroup(*instance)
	if len(flattenedGroups) > 0 {
		placementGroupMap = flattenedGroups[0]
		// Inherit compliant_only if it already exists in state
		if compliantOnly, ok := d.GetOk("placement_group.0.compliant_only"); ok {
			placementGroupMap["compliant_only"] = compliantOnly.(bool)
		}
		d.Set("placement_group", []map[string]interface{}{placementGroupMap})
	}

	disks, swapSize := flattenInstanceDisks(instanceDisks)
	d.Set("disk", disks)
	d.Set("swap_size", swapSize)

	diskLabelIDMap := make(map[int]string, len(instanceDisks))
	for _, disk := range instanceDisks {
		diskLabelIDMap[disk.ID] = disk.Label
	}

	configs := flattenInstanceConfigs(instanceConfigs, diskLabelIDMap)

	d.Set("config", configs)
	if len(instanceConfigs) == 1 {
		defaultConfig := instanceConfigs[0]

		if _, ok := d.GetOk("interface"); ok {
			flattenedInterfaces := helper.FlattenInterfaces(defaultConfig.Interfaces)
			d.Set("interface", flattenedInterfaces)
		}

		d.Set("boot_config_label", defaultConfig.Label)
	}

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_instance")

	client := meta.(*helper.ProviderMeta).Client

	if err := validateBooted(ctx, d); err != nil {
		return diag.Errorf("failed to validate: %v", err)
	}

	bootConfig := 0
	createOpts := linodego.InstanceCreateOptions{
		Region:         d.Get("region").(string),
		Type:           d.Get("type").(string),
		Label:          d.Get("label").(string),
		Group:          d.Get("group").(string),
		BackupsEnabled: d.Get("backups_enabled").(bool),
		PrivateIP:      d.Get("private_ip").(bool),
		DiskEncryption: linodego.InstanceDiskEncryption(
			d.Get("disk_encryption").(string),
		),
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	if firewallID, ok := d.GetOk("firewall_id"); ok {
		createOpts.FirewallID = firewallID.(int)
	}

	if interfaces, interfacesOk := d.GetOk("interface"); interfacesOk {
		interfaces := interfaces.([]interface{})

		createOpts.Interfaces = make([]linodego.InstanceConfigInterfaceCreateOptions, len(interfaces))

		for i, ni := range interfaces {
			createOpts.Interfaces[i] = helper.ExpandConfigInterface(ni.(map[string]interface{}))
		}
	}

	if _, metadataOk := d.GetOk("metadata.0"); metadataOk {
		var metadata linodego.InstanceMetadataOptions

		if userData, userDataOk := d.GetOk("metadata.0.user_data"); userDataOk {
			metadata.UserData = userData.(string)
		}

		createOpts.Metadata = &metadata
	}

	createOpts.PlacementGroup = getPlacementGroupCreateOptions(ctx, d)

	_, disksOk := d.GetOk("disk")
	_, configsOk := d.GetOk("config")
	bootedNull := d.GetRawConfig().GetAttr("booted").IsNull()
	booted := d.Get("booted").(bool)

	// If we don't have disks and we don't have configs, use the single API call approach
	if !disksOk && !configsOk {
		for _, key := range d.Get("authorized_keys").([]interface{}) {
			if key == nil {
				return diag.Errorf("invalid input for authorized_keys: keys cannot be empty or null")
			}

			createOpts.AuthorizedKeys = append(createOpts.AuthorizedKeys, key.(string))
		}
		for _, user := range d.Get("authorized_users").([]interface{}) {
			if user == nil {
				return diag.Errorf("invalid input for authorized_users: users cannot be empty or null")
			}

			createOpts.AuthorizedUsers = append(createOpts.AuthorizedUsers, user.(string))
		}
		createOpts.RootPass = d.Get("root_pass").(string)
		if createOpts.RootPass == "" {
			var err error
			createOpts.RootPass, err = helper.CreateRandomRootPassword()
			if err != nil {
				return diag.FromErr(err)
			}
		}

		createOpts.Image = d.Get("image").(string)

		createOpts.Booted = &boolTrue

		if !d.GetRawConfig().GetAttr("booted").IsNull() {
			createOpts.Booted = &booted
		}

		createOpts.BackupID = d.Get("backup_id").(int)
		if swapSize := d.Get("swap_size").(int); swapSize > 0 {
			createOpts.SwapSize = &swapSize
		}

		createOpts.StackScriptID = d.Get("stackscript_id").(int)

		if stackscriptDataRaw, ok := d.GetOk("stackscript_data"); ok {
			stackscriptData, ok := stackscriptDataRaw.(map[string]interface{})
			if !ok {
				return diag.Errorf("Error parsing stackscript_data: expected map[string]interface{}")
			}
			createOpts.StackScriptData = make(map[string]string, len(stackscriptData))
			for name, value := range stackscriptData {
				createOpts.StackScriptData[name] = value.(string)
			}
		}
	} else {
		createOpts.Booted = &boolFalse // necessary to prepare disks and configs
	}

	createPoller, err := client.NewEventPollerWithoutEntity(linodego.EntityLinode, linodego.ActionLinodeCreate)
	if err != nil {
		return diag.Errorf("failed to initialize event poller: %s", err)
	}

	tflog.Debug(ctx, "client.CreateInstance(...)", map[string]any{
		"options": createOpts,
	})

	instance, err := client.CreateInstance(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Instance: %s", err)
	}

	ctx = tflog.SetField(ctx, "id", instance.ID)

	d.SetId(fmt.Sprintf("%d", instance.ID))
	createPoller.EntityID = instance.ID

	var ips []string
	for _, ip := range instance.IPv4 {
		ips = append(ips, ip.String())
	}

	d.Set("ipv4", ips)
	d.Set("ipv6", instance.IPv6)

	for _, address := range instance.IPv4 {
		if private := privateIP(*address); private {
			d.Set("private_ip_address", address.String())
		} else {
			d.Set("ip_address", address.String())
		}
	}

	updateOpts := linodego.InstanceUpdateOptions{}
	doUpdate := false

	watchdogEnabled := d.Get("watchdog_enabled").(bool)
	if !watchdogEnabled {
		doUpdate = true
		updateOpts.WatchdogEnabled = &watchdogEnabled
	}

	if _, alertsOk := d.GetOk("alerts.0"); alertsOk {
		doUpdate = true
		updateOpts.Alerts = &linodego.InstanceAlert{}

		// TODO(displague) only set specified alerts
		updateOpts.Alerts.CPU = d.Get("alerts.0.cpu").(int)
		updateOpts.Alerts.IO = d.Get("alerts.0.io").(int)
		updateOpts.Alerts.NetworkIn = d.Get("alerts.0.network_in").(int)
		updateOpts.Alerts.NetworkOut = d.Get("alerts.0.network_out").(int)
		updateOpts.Alerts.TransferQuota = d.Get("alerts.0.transfer_quota").(int)
	}

	if doUpdate {
		ctx = populateLogAttributes(ctx, d)
		tflog.Debug(ctx, "client.UpdateInstance(...)", map[string]any{
			"options": updateOpts,
		})

		instance, err = client.UpdateInstance(ctx, instance.ID, updateOpts)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Look up tables for any disks and configs we create
	// - so configs and initrd can reference disks by label
	// - so configs can be referenced as a boot_config_label param
	var diskIDLabelMap map[string]int
	var configIDLabelMap map[string]int

	if disksOk {
		tflog.Debug(ctx, "Waiting for instance creation to complete before provisioning disks")

		_, err = createPoller.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
		if err != nil {
			return diag.Errorf("Error waiting for Instance to finish creating: %s", err)
		}

		tflog.Debug(ctx, "Instance is ready, provisioning disks")

		diskSpecs := d.Get("disk").([]interface{})
		diskIDLabelMap = make(map[string]int, len(diskSpecs))

		for _, diskSpec := range diskSpecs {
			diskSpec := diskSpec.(map[string]interface{})

			instanceDisk, err := createInstanceDisk(ctx, client, *instance, diskSpec, d)
			if err != nil {
				return diag.FromErr(err)
			}
			diskIDLabelMap[instanceDisk.Label] = instanceDisk.ID
		}
	}

	if configsOk {
		configSpecs := d.Get("config").([]interface{})
		detacher := makeVolumeDetacher(client, d)

		configIDMap, err := createInstanceConfigsFromSet(ctx, client, instance.ID, configSpecs, diskIDLabelMap, detacher)
		if err != nil {
			return diag.FromErr(err)
		}

		configIDLabelMap = make(map[string]int, len(configIDMap))
		for k, v := range configIDMap {
			if len(configIDLabelMap) == 1 {
				bootConfig = k
			}

			configIDLabelMap[v.Label] = k
		}
		bootConfigLabel := d.Get("boot_config_label").(string)
		if bootConfigLabel != "" {
			if foundConfig, found := configIDLabelMap[bootConfigLabel]; found {
				bootConfig = foundConfig
			} else {
				return diag.Errorf("Error setting boot_config_label: Config label '%s' not found", bootConfigLabel)
			}
		}
	}

	if ipv4Shared, ok := d.GetOk("shared_ipv4"); ok {
		shareOpts := linodego.IPAddressesShareOptions{
			IPs:      helper.ExpandStringSet(ipv4Shared.(*schema.Set)),
			LinodeID: instance.ID,
		}

		tflog.Debug(ctx, "client.ShareIPAddresses(...)", map[string]any{
			"options": shareOpts,
		})

		err = client.ShareIPAddresses(ctx, shareOpts)
		if err != nil {
			return diag.Errorf("failed to share ipv4 addresses with instance: %s", err)
		}
	}

	targetStatus := linodego.InstanceRunning

	if createOpts.Booted == nil || !*createOpts.Booted {
		if disksOk && configsOk && (bootedNull || booted) {
			p, err := client.NewEventPoller(ctx, instance.ID, linodego.EntityLinode, linodego.ActionLinodeBoot)
			if err != nil {
				return diag.Errorf("failed to initialize event poller: %s", err)
			}

			tflog.Debug(ctx, "client.BootInstance(...)", map[string]any{
				"config_id": bootConfig,
			})

			if err = client.BootInstance(ctx, instance.ID, bootConfig); err != nil {
				return diag.Errorf("Error booting Linode instance %d: %s", instance.ID, err)
			}

			event, err := p.WaitForFinished(
				ctx, getDeadlineSeconds(ctx, d),
			)
			if err != nil {
				return diag.Errorf("Error booting Linode instance %d: %s", instance.ID, err)
			}

			tflog.Debug(ctx, "Instance finished booting", map[string]any{
				"event_id": event.ID,
			})
		} else {
			targetStatus = linodego.InstanceOffline
		}
	}

	// If the instance has implicit disks and config with no specified image it will not boot.
	if !(disksOk && configsOk) && len(createOpts.Image) < 1 {
		targetStatus = linodego.InstanceOffline
	}

	if !meta.(*helper.ProviderMeta).Config.SkipInstanceReadyPoll {
		tflog.Debug(ctx, "Waiting for instance to reach target status", map[string]any{
			"target_status": targetStatus,
		})

		if _, err = client.WaitForInstanceStatus(ctx, instance.ID, targetStatus, getDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("timed-out waiting for Linode instance %d to reach status %s: %s", instance.ID, targetStatus, err)
		}
	}

	return readResource(ctx, d, meta)
}

func findDiskByFS(disks []linodego.InstanceDisk, fs linodego.DiskFilesystem) *linodego.InstanceDisk {
	for _, disk := range disks {
		if disk.Filesystem == fs {
			return &disk
		}
	}
	return nil
}

// adjustSwapSizeIfNeeded handles changes to the swap_size attribute if needed. If there is a change, this means
// resizing the underlying main/swap disks on the instance to match the declared swap size allocation.
//
// returns bool describing whether the linode needs to be restarted.
func adjustSwapSizeIfNeeded(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client, instance *linodego.Instance,
) (bool, error) {
	if !d.HasChange("swap_size") {
		return false, nil
	}

	// If the swap_size attribute is set, there are two default disks attached to the instance (the main disk of type ext4
	// and a swap disk), as custom disk configuration via "disk" nested attributes conflicts with the swap_size.
	bootDisk, swapDisk, err := getInstanceDefaultDisks(ctx, instance.ID, client)
	if err != nil {
		return false, err
	}

	oldSwapVal, newSwapVal := d.GetChange("swap_size")
	oldSwap, newSwap := oldSwapVal.(int), newSwapVal.(int)
	diff := newSwap - oldSwap
	newBootDiskSize := bootDisk.Size - diff

	toResize := []struct {
		size int
		disk *linodego.InstanceDisk
	}{
		{
			size: newBootDiskSize,
			disk: bootDisk,
		},
		{
			size: newSwap,
			disk: swapDisk,
		},
	}

	if bootDisk.Size < newBootDiskSize {
		// swap disk needs to be downsized first to upsize main disk
		toResize[0], toResize[1] = toResize[1], toResize[0]
	}

	for _, resizeOp := range toResize {
		if err := changeInstanceDiskSize(ctx, client, *instance, *resizeOp.disk, resizeOp.size, d); err != nil {
			return true, err
		}
	}
	return true, nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_instance")

	client := meta.(*helper.ProviderMeta).Client
	skipImplicitReboots := meta.(*helper.ProviderMeta).Config.SkipImplicitReboots
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	if err := validateBooted(ctx, d); err != nil {
		return diag.Errorf("failed to validate: %v", err)
	}

	instance, err := client.GetInstance(ctx, id)
	if err != nil {
		return diag.Errorf("Error fetching data about the current linode: %s", err)
	}

	updateOpts := linodego.InstanceUpdateOptions{}
	simpleUpdate := false

	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)
		simpleUpdate = true
	}
	if d.HasChange("group") {
		newGroup := d.Get("group").(string)
		updateOpts.Group = &newGroup
		simpleUpdate = true
	}
	if d.HasChange("tags") {
		tags := helper.ExpandStringSet(d.Get("tags").(*schema.Set))
		updateOpts.Tags = &tags
		simpleUpdate = true
	}
	if d.HasChange("watchdog_enabled") {
		watchdogEnabled := d.Get("watchdog_enabled").(bool)
		updateOpts.WatchdogEnabled = &watchdogEnabled
		simpleUpdate = true
	}
	if d.HasChange("alerts") {
		updateOpts.Alerts = &linodego.InstanceAlert{}
		updateOpts.Alerts.CPU = d.Get("alerts.0.cpu").(int)
		updateOpts.Alerts.IO = d.Get("alerts.0.io").(int)
		updateOpts.Alerts.NetworkIn = d.Get("alerts.0.network_in").(int)
		updateOpts.Alerts.NetworkOut = d.Get("alerts.0.network_out").(int)
		updateOpts.Alerts.TransferQuota = d.Get("alerts.0.transfer_quota").(int)
		simpleUpdate = true
	}

	if d.HasChange("placement_group.0.id") {
		oldPGID, newPGID := d.GetChange("placement_group.0.id")

		tflog.Debug(ctx, "Attempting Placement Group reassignment", map[string]any{
			"old_pg_id": oldPGID,
			"new_pg_id": newPGID,
		})

		if err := reassignPlacementGroup(
			ctx,
			client,
			instance.ID,
			oldPGID.(int),
			newPGID.(int),
			helper.SDKv2UnwrapOptionalConfigAttr[bool](ctx, d, "placement_group.0.compliant_only"),
		); err != nil {
			return diag.Errorf("failed to update Linode Placement Group: %s", err)
		}
	}

	// apply staged simple updates early
	if simpleUpdate {
		instanceID := instance.ID

		tflog.Debug(ctx, "client.UpdateInstance(...)", map[string]any{
			"options": updateOpts,
		})

		if instance, err = client.UpdateInstance(ctx, instance.ID, updateOpts); err != nil {
			return diag.Errorf("Error updating Instance %d: %s", instanceID, err)
		}
	}

	if d.HasChange("backups_enabled") {
		backupsEnabled := d.Get("backups_enabled").(bool)

		tflog.Info(ctx, "Updating backups enrollment", map[string]any{
			"backups_enabled": backupsEnabled,
		})

		if backupsEnabled {
			tflog.Debug(ctx, "client.EnableInstanceBackups(...)")
			if err = client.EnableInstanceBackups(ctx, instance.ID); err != nil {
				return diag.FromErr(err)
			}
		} else {
			tflog.Debug(ctx, "client.CancelInstanceBackups(...)")
			if err = client.CancelInstanceBackups(ctx, instance.ID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	rebootInstance := false

	if d.HasChange("private_ip") {
		if _, ok := d.GetOk("private_ip"); !ok {
			return diag.Errorf("Error removing private IP address for Instance %d: Removing a Private IP "+
				"address must be handled through a support ticket", instance.ID)
		}

		tflog.Info(ctx, "client.AddInstanceIPAddress(...)", map[string]any{
			"public": false,
		})

		privateIP, err := client.AddInstanceIPAddress(ctx, instance.ID, false)
		if err != nil {
			return diag.Errorf("Error activating private networking on Instance %d: %s", instance.ID, err)
		}
		d.Set("private_ip_address", privateIP.Address)
		rebootInstance = true
	}

	// If the region has changed,
	// we should migrate the Linode.
	if d.HasChange("region") {
		instance, err = applyInstanceMigration(
			ctx,
			d,
			&client,
			instance,
			d.Get("region").(string),
		)
		if err != nil {
			return diag.Errorf("failed to migrate instance: %s", err)
		}
	}

	oldSpec, newSpec, err := getInstanceTypeChange(ctx, d, &client)
	if err != nil {
		return diag.Errorf("Error getting resize info for instance: %s", err)
	}
	upsized := newSpec.Disk > oldSpec.Disk

	if upsized {
		// The linode was upsized; apply before disk changes to allocate more disk
		if instance, err = applyInstanceTypeChange(ctx, d, &client, instance, newSpec); err != nil {
			return diag.Errorf("failed to change instance type: %s", err)
		}
		rebootInstance = true
	}

	// We only need to do this if explicit disks are defined
	if d.GetRawConfig().GetAttr("image").IsNull() {
		if didChange, err := applyInstanceDiskSpec(ctx, d, &client, instance, newSpec); err == nil && didChange {
			rebootInstance = true
		} else if err != nil && newSpec.Disk < oldSpec.Disk && !d.HasChange("disk") {
			// Linode was downsized but the pre-existing disk config does not fit new instance spec
			// This might mean the user tried to downsize an instance with an implicit, default
			return diag.Errorf("failed to apply instance disk spec: %s."+downsizeFailedMessage, err)
		} else if err != nil {
			return diag.Errorf("failed to apply instance disk spec: %s", err)
		}
	}

	if oldSpec.ID != newSpec.ID && !upsized {
		// linode was downsized or changed to a type with the same disk allocation
		if instance, err = applyInstanceTypeChange(ctx, d, &client, instance, newSpec); err != nil {
			return diag.Errorf("failed to change instance type: %s", err)
		}
	}

	if didChange, err := adjustSwapSizeIfNeeded(ctx, d, &client, instance); err != nil {
		return diag.FromErr(err)
	} else if didChange {
		rebootInstance = true
	}

	diskIDLabelMap, err := getInstanceDiskLabelIDMap(ctx, client, d, instance.ID)
	if err != nil {
		return diag.Errorf("failed to get disk label to ID mappings")
	}

	bootConfig := 0
	bootConfigLabel := d.Get("boot_config_label").(string)

	tfConfigsOld, tfConfigsNew := d.GetChange("config")
	didChangeConfig, updatedConfigMap, updatedConfigs, err := updateInstanceConfigs(
		ctx, client, d, *instance, tfConfigsOld, tfConfigsNew, diskIDLabelMap, bootConfigLabel)
	if err != nil {
		return diag.FromErr(err)
	}
	rebootInstance = rebootInstance || didChangeConfig

	if bootConfigLabel != "" {
		if foundConfig, found := updatedConfigMap[bootConfigLabel]; found {
			bootConfig = foundConfig
		} else {
			return diag.Errorf("Error setting boot_config_label: Config label '%s' not found", bootConfigLabel)
		}
	} else if len(updatedConfigs) > 0 {
		bootConfig = updatedConfigs[0].ID
	}

	booted := d.Get("booted").(bool)
	bootedNull := d.GetRawConfig().GetAttr("booted").IsNull()

	if d.HasChange("interface") {
		interfaces := d.Get("interface").([]interface{})

		expandedInterfaces := helper.ExpandConfigInterfaces(ctx, interfaces)
		config, err := client.GetInstanceConfig(ctx, id, bootConfig)
		if err != nil {
			return diag.Errorf("failed to get config %d: %s", bootConfig, err)
		}

		powerOffRequired := VPCInterfaceIncluded(config.Interfaces, expandedInterfaces)

		tflog.Debug(ctx, "Updating instance config for interface changes", map[string]any{
			"config_id": bootConfig,
		})

		instance, err := client.GetInstance(ctx, id)
		if err != nil {
			return diag.Errorf("Error fetching data about the current linode: %s", err)
		}

		// we should power on Linode after updating of the interfaces if
		// it's currently on and booted attribute is unset by the user.
		// Otherwise, it will stay off (if it's already off) or be handled by
		// `handleBootedUpdate` (if booted is set to an explicit value)
		shouldPowerOn := bootedNull && powerOffRequired && instance.Status == linodego.InstanceRunning

		if powerOffRequired {
			if err := ShutdownInstanceForVPCInterfaceUpdate(
				ctx, &client, skipImplicitReboots, id, helper.GetDeadlineSeconds(ctx, d),
			); err != nil {
				return diag.FromErr(err)
			}
		}

		// reboot won't be needed if we power off the Linode during update
		rebootInstance = !powerOffRequired
		configUpdateOpts := linodego.InstanceConfigUpdateOptions{
			Interfaces: expandedInterfaces,
		}

		tflog.Debug(
			ctx,
			"client.UpdateInstanceConfig(...)",
			map[string]any{"options": configUpdateOpts},
		)

		if _, err := client.UpdateInstanceConfig(
			ctx, instance.ID, bootConfig, configUpdateOpts,
		); err != nil {
			return diag.Errorf("failed to set boot config interfaces: %s", err)
		}

		if shouldPowerOn {
			if diag := BootInstanceAfterVPCInterfaceUpdate(
				ctx, meta.(*helper.ProviderMeta), id, bootConfig, helper.GetDeadlineSeconds(ctx, d),
			); diag != nil {
				return diag
			}
		}

	}

	if d.HasChange("shared_ipv4") {
		sharedIPs := helper.ExpandStringSet(d.Get("shared_ipv4").(*schema.Set))

		shareOpts := linodego.IPAddressesShareOptions{
			IPs:      sharedIPs,
			LinodeID: instance.ID,
		}

		tflog.Debug(ctx, "client.ShareIPAddresses(...)", map[string]any{
			"options": shareOpts,
		})

		err = client.ShareIPAddresses(ctx, shareOpts)
		if err != nil {
			return diag.Errorf("failed to share ipv4 addresses with instance: %s", err)
		}
	}

	// Don't reboot if the Linode should be powered off
	if !bootedNull && !booted {
		rebootInstance = false
	}

	// Only reboot the instance if implicit reboots are not skipped
	if skipImplicitReboots {
		rebootInstance = false
	}

	if rebootInstance && len(diskIDLabelMap) > 0 && len(updatedConfigMap) > 0 && bootConfig > 0 {
		ctx := tflog.SetField(ctx, "config_id", bootConfig)

		tflog.Info(ctx, "Implicitly rebooting instance")

		p, err := client.NewEventPoller(ctx, id, linodego.EntityLinode, linodego.ActionLinodeReboot)
		if err != nil {
			return diag.Errorf("failed to initialize event poller: %s", err)
		}

		tflog.Debug(ctx, "client.RebootInstance(...)")

		err = client.RebootInstance(ctx, instance.ID, bootConfig)
		if err != nil {
			return diag.Errorf("Error rebooting Instance %d: %s", instance.ID, err)
		}

		tflog.Debug(ctx, "Waiting for instance reboot to complete")

		_, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d))
		if err != nil {
			return diag.Errorf("Error waiting for Instance %d to finish rebooting: %s", instance.ID, err)
		}

		tflog.Debug(ctx, "Instance has finished rebooting")

		if _, err = client.WaitForInstanceStatus(
			ctx, instance.ID, linodego.InstanceRunning, getDeadlineSeconds(ctx, d),
		); err != nil {
			return diag.Errorf("Timed-out waiting for Linode instance %d to boot: %s", instance.ID, err)
		}
	}

	if err := handleBootedUpdate(ctx, d, meta, instance.ID, bootConfig); err != nil {
		return diag.Errorf("failed to handle booted update: %s", err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_instance")

	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int", d.Id())
	}

	p, err := client.NewEventPoller(ctx, id, linodego.EntityLinode, linodego.ActionLinodeDelete)
	if err != nil {
		return diag.Errorf("failed to initialize event poller: %s", err)
	}

	tflog.Debug(ctx, "client.DeleteInstance(...)")

	err = client.DeleteInstance(ctx, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode instance %d: %s", id, err)
	}

	if !meta.(*helper.ProviderMeta).Config.SkipInstanceDeletePoll {
		// Wait for full deletion to assure volumes are detached
		if _, err = p.WaitForFinished(ctx, getDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("failed to wait for instance %d to be deleted: %s", id, err)
		}
	}

	d.SetId("")
	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": d.Id(),
	})
}
