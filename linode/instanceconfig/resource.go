package instanceconfig

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	tflog.Debug(ctx, "Import linode_instance_config", map[string]any{
		"id": d.Id(),
	})

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		// Validate that this is an ID by making sure it can be converted into an int
		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid config ID: %v", err)
		}

		instID, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid instance ID: %v", err)
		}

		d.SetId(s[1])
		d.Set("linode_id", instID)
	}

	err := readResource(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import %v as instance config: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_instance_config")

	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	cfg, err := client.GetInstanceConfig(ctx, linodeID, id)
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

	// We want to guarantee that we're resolving a public IPv4 address
	instNetworking, err := client.GetInstanceIPAddresses(ctx, linodeID)
	if err != nil {
		return diag.Errorf("failed to get instance networking: %s", err)
	}

	configBooted, err := isConfigBooted(ctx, &client, inst, cfg.ID)
	if err != nil {
		return diag.Errorf("failed to check instance boot status: %s", err)
	}

	d.Set("linode_id", linodeID)
	d.Set("label", cfg.Label)
	d.Set("comments", cfg.Comments)
	d.Set("kernel", cfg.Kernel)
	d.Set("memory_limit", cfg.MemoryLimit)
	d.Set("root_device", cfg.RootDevice)
	d.Set("run_level", cfg.RunLevel)
	d.Set("virt_mode", cfg.VirtMode)
	d.Set("interface", flattenInterfaces(cfg.Interfaces))
	d.Set("booted", configBooted)

	if cfg.Devices != nil {
		d.Set("devices", flattenDeviceMapToNamedBlock(*cfg.Devices))
		d.Set("device", flattenDeviceMapToBlock(*cfg.Devices))
	}

	if cfg.Helpers != nil {
		d.Set("helpers", flattenHelpers(*cfg.Helpers))
	}

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": instNetworking.IPv4.Public[0].Address,
	})

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_instance_config")

	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	inst, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	createOpts := linodego.InstanceConfigCreateOptions{
		Label:       d.Get("label").(string),
		Comments:    d.Get("comments").(string),
		Helpers:     expandHelpers(d.Get("helpers")),
		Interfaces:  expandInterfaces(ctx, d.Get("interface").([]any)),
		MemoryLimit: d.Get("memory_limit").(int),
		Kernel:      d.Get("kernel").(string),
		RunLevel:    d.Get("run_level").(string),
		VirtMode:    d.Get("virt_mode").(string),
	}

	if rootDevice, ok := d.GetOk("root_device"); ok {
		rootDeviceStr := rootDevice.(string)
		createOpts.RootDevice = &rootDeviceStr
	}

	var devices *linodego.InstanceConfigDeviceMap
	if devicesBlock, ok := d.GetOk("device"); ok {
		devices = expandDevicesBlock(devicesBlock)
	} else if devicesBlock, ok := d.GetOk("devices"); ok {
		devices = expandDevicesNamedBlock(devicesBlock)
	}
	if devices != nil {
		createOpts.Devices = *devices
	}

	tflog.Debug(ctx, "Sending config creation API request", map[string]any{
		"body": createOpts,
	})

	cfg, err := client.CreateInstanceConfig(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance config: %s", err)
	}

	ctx = tflog.SetField(ctx, "config_id", cfg.ID)

	d.SetId(strconv.Itoa(cfg.ID))

	if !d.GetRawConfig().GetAttr("booted").IsNull() {
		if err := applyBootStatus(ctx, &client, inst, cfg.ID, helper.GetDeadlineSeconds(ctx, d),
			d.Get("booted").(bool)); err != nil {
			return diag.Errorf("failed to update boot status: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_instance_config")

	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"id":        id,
		"linode_id": linodeID,
	})

	tflog.Debug(ctx, "Update resource")

	inst, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	putRequest := linodego.InstanceConfigUpdateOptions{}
	shouldUpdate := false

	if d.HasChange("comments") {
		putRequest.Comments = d.Get("comments").(string)
		shouldUpdate = true
	}

	if d.HasChange("device") {
		if devices, ok := d.GetOk("device"); ok {
			putRequest.Devices = expandDevicesBlock(devices)
		}
		shouldUpdate = true
	}

	if d.HasChange("devices") {
		if devices, ok := d.GetOk("devices"); ok {
			putRequest.Devices = expandDevicesNamedBlock(devices)
		}
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

	if d.HasChange("interface") {
		putRequest.Interfaces = expandInterfaces(ctx, d.Get("interface").([]any))
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Debug(ctx, "Update detected, sending config PUT request to API", map[string]any{
			"body": putRequest,
		})
		if _, err := client.UpdateInstanceConfig(ctx, linodeID, int(id), putRequest); err != nil {
			return diag.Errorf("failed to update instance config: %s", err)
		}
	}

	// We should not use `HasChange(...)` here because of possible mid-apply changes
	if !d.GetRawConfig().GetAttr("booted").IsNull() {
		if err := applyBootStatus(ctx, &client, inst, id, helper.GetDeadlineSeconds(ctx, d),
			d.Get("booted").(bool)); err != nil {
			return diag.Errorf("failed to update boot status: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_instance_config")

	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse id: %s", err)
	}

	linodeID := d.Get("linode_id").(int)

	inst, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	// Shutdown the instance if the config is in use
	if booted, err := isConfigBooted(ctx, &client, inst, id); err != nil {
		return diag.Errorf("failed to check if config is booted: %s", err)
	} else if booted {
		tflog.Info(ctx, "Shutting down instance for config deletion")

		p, err := client.NewEventPoller(ctx, inst.ID, linodego.EntityLinode, linodego.ActionLinodeShutdown)
		if err != nil {
			return diag.Errorf("failed to poll for events: %s", err)
		}

		if err := client.ShutdownInstance(ctx, inst.ID); err != nil {
			return diag.Errorf("failed to shutdown instance: %s", err)
		}

		if _, err := p.WaitForFinished(ctx, helper.GetDeadlineSeconds(ctx, d)); err != nil {
			return diag.Errorf("failed to wait for instance shutdown: %s", err)
		}
		tflog.Debug(ctx, "Instance shutdown complete")
	}

	tflog.Debug(ctx, "Deleting instance config")

	err = client.DeleteInstanceConfig(ctx, linodeID, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance Config %d: %s", id, err)
	}
	d.SetId("")

	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": d.Get("linode_id").(int),
		"id":        d.Id(),
	})
}
