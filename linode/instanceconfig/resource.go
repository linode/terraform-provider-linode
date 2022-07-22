package instanceconfig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"log"
	"strconv"
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

	configBooted, err := isConfigBooted(ctx, &client, inst, cfg.ID)
	if err != nil {
		return diag.Errorf("failed to check instance boot status: %s", err)
	}

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

	inst, err := client.GetInstance(ctx, linodeID)
	if err != nil {
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	createOpts := linodego.InstanceConfigCreateOptions{
		Label:       d.Get("label").(string),
		Comments:    d.Get("comments").(string),
		Helpers:     expandHelpers(d.Get("helpers")),
		Interfaces:  expandInterfaces(d.Get("interface").([]any)),
		MemoryLimit: d.Get("memory_limit").(int),
		Kernel:      d.Get("kernel").(string),
		RunLevel:    d.Get("run_level").(string),
		VirtMode:    d.Get("virt_mode").(string),
	}

	if rootDevice, ok := d.GetOk("root_device"); ok {
		rootDeviceStr := rootDevice.(string)
		createOpts.RootDevice = &rootDeviceStr
	}

	if devices, ok := d.GetOk("devices"); ok {
		createOpts.Devices = *expandDeviceMap(devices)
	}

	cfg, err := client.CreateInstanceConfig(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance config: %s", err)
	}

	d.SetId(strconv.Itoa(cfg.ID))

	if !d.GetRawConfig().GetAttr("booted").IsNull() {
		if err := applyBootStatus(ctx, &client, inst, cfg.ID, helper.GetDeadlineSeconds(ctx, d), d.Get("booted").(bool)); err != nil {
			return diag.Errorf("failed to update boot status: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	linodeID := d.Get("linode_id").(int)

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

	if d.HasChange("interface") {
		putRequest.Interfaces = expandInterfaces(d.Get("interface").([]any))
	}

	if shouldUpdate {
		if _, err := client.UpdateInstanceConfig(ctx, linodeID, int(id), putRequest); err != nil {
			return diag.Errorf("failed to update instance config: %s", err)
		}
	}

	if !d.GetRawConfig().GetAttr("booted").IsNull() && d.HasChange("booted") {
		if err := applyBootStatus(ctx, &client, inst, int(id), helper.GetDeadlineSeconds(ctx, d), d.Get("booted").(bool)); err != nil {
			return diag.Errorf("failed to update boot status: %s", err)
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
