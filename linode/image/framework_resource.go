package image

import (
	"context"
	"crypto/md5" // #nosec G501 -- endpoint expecting md5
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	DefaultVolumeCreateTimeout = 30 * time.Minute
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_image",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
				TimeoutOpts: &timeouts.Opts{
					Create: true,
				},
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func createResourceFromUpload(
	ctx context.Context, plan *ResourceModel, client *linodego.Client, resp *resource.CreateResponse, timeoutSeconds int,
) *linodego.Image {
	tflog.Debug(ctx, "Create linode_image from file uploading")

	imageReader := openImageFile(plan.FilePath.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return nil
	}

	defer func() {
		if err := imageReader.Close(); err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to close image reader: %q\n", err.Error()))
		}
	}()

	createOpts := linodego.ImageCreateUploadOptions{
		Region:      plan.Region.ValueString(),
		Label:       plan.Label.ValueString(),
		Description: plan.Description.ValueString(),
		CloudInit:   plan.CloudInit.ValueBool(),
	}

	tflog.Trace(ctx, "client.CreateImageUpload(...)", map[string]any{
		"options": createOpts,
	})

	image, uploadURL, err := client.CreateImageUpload(ctx, createOpts)
	if image != nil && len(image.ID) > 0 {
		addImageResource(ctx, resp, image.ID)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to Create Image (%s) from File Uploading",
				plan.Label.ValueString(),
			),
			err.Error(),
		)
		return image
	}

	ctx = tflog.SetField(ctx, "image_id", image.ID)

	uploadImageAndStoreHash(ctx, plan, client, uploadURL, imageReader, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return image
	}

	tflog.Debug(ctx, "Waiting for a single image to be ready")
	tflog.Trace(ctx, "client.WaitForImageStatus(...)", map[string]any{
		"status": linodego.ImageStatusAvailable,
	})

	image, err = client.WaitForImageStatus(
		ctx,
		image.ID,
		linodego.ImageStatusAvailable,
		timeoutSeconds,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Wait for Image to be Available", err.Error())
		return image
	}

	return refreshImage(ctx, image, client, &resp.Diagnostics)
}

func createResourceFromLinode(
	ctx context.Context, plan *ResourceModel, client *linodego.Client, resp *resource.CreateResponse, timeoutSeconds int,
) *linodego.Image {
	tflog.Debug(ctx, "Create linode_image from a Linode instance")

	linodeID := helper.FrameworkSafeInt64ToInt(plan.LinodeID.ValueInt64(), &resp.Diagnostics)
	diskID := helper.FrameworkSafeInt64ToInt(plan.DiskID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return nil
	}

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, timeoutSeconds,
	); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to wait for Linode Instance %d Disk %d to become ready for taking an Image",
				linodeID, diskID,
			),
			err.Error(),
		)
		return nil
	}

	createOpts := linodego.ImageCreateOptions{
		DiskID:      diskID,
		Label:       plan.Label.ValueString(),
		Description: plan.Description.ValueString(),
		CloudInit:   plan.CloudInit.ValueBool(),
	}
	tflog.Trace(ctx, "client.CreateImage(...)", map[string]any{
		"options": createOpts,
	})

	image, err := client.CreateImage(ctx, createOpts)
	if image != nil && len(image.ID) > 0 {
		addImageResource(ctx, resp, image.ID)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating a Linode Image", err.Error())
		return image
	}

	ctx = populateLogAttributes(ctx, image.ID)
	tflog.Debug(ctx, "Waiting for a single image to be ready")

	tflog.Trace(ctx, "client.WaitForInstanceDiskStatus(...)", map[string]any{
		"status": "ready",
	})

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, timeoutSeconds,
	); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to wait for linode instance %d disk %d to become ready while taking an image",
				linodeID, diskID,
			),
			err.Error(),
		)
	}

	// Override unknown hash to null.
	// Hash is only known when uploading image from file
	plan.FileHash = helper.KeepOrUpdateValue(plan.FileHash, types.StringNull(), true)

	return refreshImage(ctx, image, client, &resp.Diagnostics)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "Create "+r.Config.Name)

	var plan ResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, DefaultVolumeCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutSeconds := helper.FrameworkSafeFloat64ToInt(createTimeout.Seconds(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var image *linodego.Image
	if !plan.LinodeID.IsNull() && plan.FilePath.IsNull() {
		image = createResourceFromLinode(ctx, &plan, client, resp, timeoutSeconds)
	} else {
		image = createResourceFromUpload(ctx, &plan, client, resp, timeoutSeconds)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	plan.FlattenImage(image, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.ID = types.StringValue(image.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	imageID := state.ID.ValueString()

	ctx = populateLogAttributes(ctx, imageID)

	image, err := client.GetImage(ctx, imageID)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Image No Longer Exists",
				fmt.Sprintf("Removing image %s from the state", imageID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to Get Image", err.Error())
		return
	}

	state.FlattenImage(image, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageID := plan.ID.ValueString()
	ctx = populateLogAttributes(ctx, imageID)

	client := r.Meta.Client

	updateOpts := linodego.ImageUpdateOptions{}
	shouldUpdate := false

	if !state.Description.Equal(plan.Description) {
		updateOpts.Description = plan.Description.ValueStringPointer()
		shouldUpdate = true
	}

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
		shouldUpdate = true
	}

	if shouldUpdate {
		tflog.Debug(ctx, "client.UpdateImage(...)", map[string]any{
			"options": updateOpts,
		})

		image, err := client.UpdateImage(ctx, imageID, updateOpts)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Update Image", err.Error())
			return
		}
		plan.FlattenImage(image, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan.CopyFrom(state, true)

	// Workaround for Crossplane issue where ID is not
	// properly populated in plan
	// See TPT-2865 for more details
	if plan.ID.ValueString() == "" {
		plan.ID = state.ID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageID := state.ID.ValueString()
	ctx = populateLogAttributes(ctx, imageID)

	client := r.Meta.Client

	tflog.Trace(ctx, "client.DeleteImage(...)")

	err := client.DeleteImage(ctx, imageID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Image", err.Error())
		return
	}
}

func populateLogAttributes(ctx context.Context, imageID string) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"image_id": imageID,
	})
}

func addImageResource(
	ctx context.Context, resp *resource.CreateResponse, id string,
) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))
}

func refreshImage(
	ctx context.Context, image *linodego.Image, client *linodego.Client, diags *diag.Diagnostics,
) *linodego.Image {
	tflog.Trace(ctx, "Enter refreshImage")
	if image == nil {
		diags.AddError(
			"Can't Refresh the Image",
			"Image is nil, and it can't be refreshed. "+
				"Please report this bug to the provider developers.",
		)
		return nil
	}

	image, err := client.GetImage(ctx, image.ID)
	if err != nil {
		diags.AddError("Can't Refresh the Image", err.Error())
	}

	return image
}

func openImageFile(imageFile string, diags *diag.Diagnostics) io.ReadCloser {
	file, err := os.Open(filepath.Clean(imageFile))
	if err != nil {
		diags.AddError(
			fmt.Sprintf("Failed to Open Image File %q", imageFile),
			err.Error(),
		)
		return nil
	}

	return file
}

func uploadImageAndStoreHash(
	ctx context.Context, plan *ResourceModel, client *linodego.Client,
	uploadURL string, image io.Reader, diags *diag.Diagnostics,
) {
	hash := md5.New() // #nosec G401 -- endpoint expecting md5
	tee := io.TeeReader(image, hash)

	tflog.Debug(ctx, "client.UploadImageToURL(...)", map[string]any{
		"upload_url": uploadURL,
	})

	if err := client.UploadImageToURL(ctx, uploadURL, tee); err != nil {
		diags.AddError(
			"Failed to Upload Image to URL",
			err.Error(),
		)
		return
	}

	plan.FileHash = types.StringValue(hex.EncodeToString(hash.Sum(nil)))
}
