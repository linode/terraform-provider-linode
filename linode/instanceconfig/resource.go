package instanceconfig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	linodeID := d.Get("linode_id").(int)

	cfg, err := client.GetInstanceConfig(ctx, linodeID, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Instance Config ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Instance Config: %s", err)
	}

	inst, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	d.Set("label", cfg.Label)
	d.Set("comments", cfg.Comments)
	d.Set("kernel", cfg.Kernel)
	d.Set("memory_limit", cfg.MemoryLimit)
	d.Set("root_device", cfg.RootDevice)
	d.Set("run_level", cfg.RunLevel)
	d.Set("virt_mode", cfg.VirtMode)
	d.Set("interface", flattenInterfaces(cfg.Interfaces))
	d.Set("booted", isInstanceBooted(inst.Status))

	if cfg.Devices != nil {
		d.Set("devices", flattenDeviceMap(*cfg.Devices))
	}

	if cfg.Helpers != nil {
		d.Set("helpers", flattenHelpers(*cfg.Helpers))
	}

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	deviceMap := expandDeviceMap(d.Get("devices"))
	if deviceMap == nil {
		return diag.Errorf("device map is expectedly nil")
	}

	createOpts := linodego.InstanceConfigCreateOptions{
		Label:       d.Get("label").(string),
		Comments:    d.Get("comments").(string),
		Devices:     *deviceMap,
		Helpers:     expandHelpers(d.Get("helpers")),
		Interfaces:  nil,
		MemoryLimit: d.Get("memory_limit").(int),
		Kernel:      d.Get("kernel").(string),
		RunLevel:    d.Get("run_level").(string),
		VirtMode:    d.Get("virt_mode").(string),
	}

	if rootDevice, ok := d.GetOk("root_device"); ok {
		rootDeviceStr := rootDevice.(string)
		createOpts.RootDevice = &rootDeviceStr
	}

	cfg, err := client.CreateInstanceConfig(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance config: %s", err)
	}

	d.SetId(strconv.Itoa(cfg.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	linodeID := d.Get("linode_id").(int)

	putRequest := linodego.InstanceConfigUpdateOptions{}
	shouldUpdate := false

	if d.HasChange("comments") {
		putRequest.Comments = d.Get("comments").(string)
		shouldUpdate = true
	}

	if d.HasChange("devices") {
		putRequest.Devices = expandDeviceMap(d.Get("devices"))
		shouldUpdate = true
	}

	if d.HasChange("helpers") {
		putRequest.Helpers = expandHelpers(d.Get("helpers"))
		shouldUpdate = true
	}

	if d.HasChange("kernel") {
		putRequest.Kernel = d.Get("kernel").(string)
		shouldUpdate = true
	}

	if d.HasChange("label") {
		putRequest.Label = d.Get("label").(string)
		shouldUpdate = true
	}

	if d.HasChange("memory_limit") {
		putRequest.MemoryLimit = d.Get("memory_limit").(int)
		shouldUpdate = true
	}

	if d.HasChange("root_device") {
		putRequest.RootDevice = d.Get("root_device").(string)
		shouldUpdate = true
	}

	if d.HasChange("run_level") {
		putRequest.RunLevel = d.Get("run_level").(string)
		shouldUpdate = true
	}

	if d.HasChange("virt_mode") {
		putRequest.VirtMode = d.Get("virt_mode").(string)
		shouldUpdate = true
	}

	if shouldUpdate {
		if _, err := client.UpdateInstanceConfig(ctx, linodeID, int(id), putRequest); err != nil {
			return diag.Errorf("failed to update instance config: %s", err)
		}
	}

	if !d.GetRawConfig().GetAttr("booted").IsNull() && d.HasChange("booted") {
		booted := d.Get("booted").(bool)

		if booted {
			if err := client.BootInstance(ctx, linodeID, int(id)); err != nil {
				return diag.Errorf("failed to boot to instance config: %s", err)
			}

			if _, err := client.WaitForEventFinished(ctx, linodeID, linodego.EntityLinode,
				linodego.ActionLinodeBoot, time.Now(), helper.GetDeadlineSeconds(ctx, d)); err != nil {
				return diag.Errorf("failed to wait for instance boot: %s", err)
			}
		} else {
			if err := client.ShutdownInstance(ctx, linodeID); err != nil {
				return diag.Errorf("failed to shutdown instance: %s", err)
			}
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id64, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance id %s as int", d.Id())
	}

	linodeID := d.Get("linode_id").(int)

	err = client.DeleteInstanceConfig(ctx, linodeID, int(id64))
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance Config %d: %s", id64, err)
	}
	d.SetId("")
	return nil
}

func flattenDeviceMap(deviceMap linodego.InstanceConfigDeviceMap) []map[string]any {
	result := make(map[string]any)

	reflectMap := reflect.ValueOf(deviceMap)

	for i := 0; i < reflectMap.NumField(); i++ {
		field := reflectMap.Field(i).Interface().(*linodego.InstanceConfigDevice)
		if field == nil {
			continue
		}

		fieldName := strings.ToLower(reflectMap.Type().Field(i).Name)

		result[fieldName] = map[string]any{
			"disk_id":   field.DiskID,
			"volume_id": field.VolumeID,
		}
	}

	return []map[string]any{result}
}

func flattenHelpers(helpers linodego.InstanceConfigHelpers) []map[string]any {
	result := make(map[string]any)

	result["devtmpfs_automount"] = helpers.DevTmpFsAutomount
	result["distro"] = helpers.Distro
	result["modules_dep"] = helpers.ModulesDep
	result["network"] = helpers.Network
	result["updatedb_disable"] = helpers.UpdateDBDisabled

	return []map[string]any{result}
}

func flattenInterfaces(interfaces []linodego.InstanceConfigInterface) []map[string]any {
	result := make([]map[string]any, len(interfaces))

	for i, iface := range interfaces {
		result[i] = map[string]any{
			"purpose":      iface.Purpose,
			"ipam_address": iface.IPAMAddress,
			"label":        iface.Label,
		}
	}

	return result
}

func expandDeviceMap(deviceMap any) *linodego.InstanceConfigDeviceMap {
	var result linodego.InstanceConfigDeviceMap
	deviceMapSlice := deviceMap.([]any)

	if len(deviceMapSlice) < 1 {
		return nil
	}

	devices := deviceMapSlice[0].(map[string]any)

	for k, v := range devices {
		currentDeviceSlice := v.([]any)
		if len(currentDeviceSlice) < 1 {
			continue
		}

		currentDevice := currentDeviceSlice[0].(map[string]any)

		device := linodego.InstanceConfigDevice{}

		if diskID, ok := currentDevice["disk_id"]; ok {
			device.DiskID = diskID.(int)
		}

		if volumeID, ok := currentDevice["volume_id"]; ok {
			device.VolumeID = volumeID.(int)
		}

		// Get the corresponding struct field and set it to the correct device
		field := reflect.Indirect(reflect.ValueOf(&result)).FieldByName(strings.ToUpper(k))
		field.Set(reflect.ValueOf(&device))
	}

	return &result
}

func expandHelpers(helpersRaw any) *linodego.InstanceConfigHelpers {
	helpersSlice := helpersRaw.([]any)

	if len(helpersSlice) < 1 {
		return nil
	}

	helpers := helpersSlice[0].(map[string]any)

	return &linodego.InstanceConfigHelpers{
		UpdateDBDisabled:  helpers["updatedb_disabled"].(bool),
		Distro:            helpers["distro"].(bool),
		ModulesDep:        helpers["modules_dep"].(bool),
		Network:           helpers["network"].(bool),
		DevTmpFsAutomount: helpers["devtmpfs_automount"].(bool),
	}
}

func isInstanceBooted(status linodego.InstanceStatus) bool {
	// For diffing purposes, transition states need to be treated as
	// booted == true. This is because these statuses will eventually
	// result in a powered on Linode.
	return status == linodego.InstanceRunning ||
		status == linodego.InstanceRebooting ||
		status == linodego.InstanceBooting
}
