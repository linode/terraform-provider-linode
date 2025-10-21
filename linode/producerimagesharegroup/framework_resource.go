package producerimagesharegroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_producer_image_share_group",
				IDType: types.Int64Type,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
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

	// Prepare list of images to include in the Image Share Group
	var imageShares []linodego.ImageShareGroupImage

	if !plan.Images.IsNull() && !plan.Images.IsUnknown() {
		var imageModels []ImageShareAttributesModel
		resp.Diagnostics.Append(plan.Images.ElementsAs(ctx, &imageModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, img := range imageModels {
			imageShares = append(imageShares, linodego.ImageShareGroupImage{
				ID:          img.ID.ValueString(),
				Label:       img.Label.ValueStringPointer(),
				Description: img.Description.ValueStringPointer(),
			})
		}
	}

	createOpts := linodego.ImageShareGroupCreateOptions{
		Label:       plan.Label.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		Images:      imageShares,
	}

	tflog.Debug(ctx, "client.CreateImageShareGroup(...)", map[string]any{
		"options": createOpts,
	})
	sg, err := client.CreateImageShareGroup(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Image Share Group.",
			err.Error(),
		)
		return
	}

	plan.FlattenImageShareGroup(sg, true)
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

	imageShareGroupID := helper.FrameworkSafeInt64ToInt(state.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the main Image Share Group details
	sg, err := client.GetImageShareGroup(ctx, imageShareGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read Image Share Group.",
			err.Error(),
		)
		return
	}

	// Retrieve the list of images in this Image Share Group
	imagesResp, err := client.ImageShareGroupListImageShareEntries(ctx, imageShareGroupID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read images for Image Share Group.",
			err.Error(),
		)
		return
	}

	// Convert images into simplified image share structs
	var imageShares []linodego.ImageShareGroupImage

	for _, img := range imagesResp {
		sourceID := ""
		if img.ImageSharing.SharedBy != nil && img.ImageSharing.SharedBy.SourceImageID != nil {
			sourceID = *img.ImageSharing.SharedBy.SourceImageID
		}

		imageShares = append(imageShares, linodego.ImageShareGroupImage{
			ID:          sourceID,
			Label:       linodego.Pointer(img.Label),
			Description: linodego.Pointer(img.Description),
		})
	}

	// Flatten the main share group
	state.FlattenImageShareGroup(sg, true)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the images list in the state
	var imageModels []ImageShareAttributesModel
	for _, img := range imageShares {
		imageModels = append(imageModels, ImageShareAttributesModel{
			ID:          types.StringValue(img.ID),
			Label:       types.StringPointerValue(img.Label),
			Description: types.StringPointerValue(img.Description),
		})
	}

	listVal, diag := types.ListValueFrom(ctx, imageShareGroupImage.Type(), imageModels)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Images = listVal

	// Normalize null/unknown -> []
	if state.Images.IsNull() || state.Images.IsUnknown() {
		state.Images = types.ListValueMust(imageShareGroupImage.Type(), []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client

	var plan ResourceModel
	var state ResourceModel

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageShareGroupID := helper.FrameworkSafeInt64ToInt(state.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh remote state to ensure accuracy
	sg, err := client.GetImageShareGroup(ctx, imageShareGroupID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read Image Share Group.",
			err.Error(),
		)
		return
	}

	imagesResp, err := client.ImageShareGroupListImageShareEntries(ctx, imageShareGroupID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read images for Image Share Group.",
			err.Error(),
		)
		return
	}

	var imageShares []imageData

	for _, img := range imagesResp {
		sourceID := ""
		if img.ImageSharing.SharedBy != nil && img.ImageSharing.SharedBy.SourceImageID != nil {
			sourceID = *img.ImageSharing.SharedBy.SourceImageID
		}

		imageShares = append(imageShares, imageData{
			PrivateID:   sourceID,
			SharedID:    img.ID,
			Label:       linodego.Pointer(img.Label),
			Description: linodego.Pointer(img.Description),
		})
	}

	// Build desired vs. actual sets
	var planImages []ImageShareAttributesModel
	if !plan.Images.IsNull() && !plan.Images.IsUnknown() {
		resp.Diagnostics.Append(plan.Images.ElementsAs(ctx, &planImages, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build lookup maps
	planMap := make(map[string]ImageShareAttributesModel)
	for _, img := range planImages {
		planMap[img.ID.ValueString()] = img
	}

	remoteMap := make(map[string]imageData)
	for _, img := range imageShares {
		remoteMap[img.PrivateID] = img
	}

	var toAdd []linodego.ImageShareGroupImage
	var toUpdate []struct {
		id   string
		opts linodego.ImageShareGroupUpdateImageOptions
	}
	var toRemove []string

	// Detect creates and updates
	for privateID, planImg := range planMap {
		remote, exists := remoteMap[privateID]
		if !exists {
			// Not found remotely, so add
			toAdd = append(toAdd, linodego.ImageShareGroupImage{
				ID:          privateID,
				Label:       planImg.Label.ValueStringPointer(),
				Description: planImg.Description.ValueStringPointer(),
			})
		} else {
			// Found remotely, check for changes
			labelChanged := planImg.Label.ValueString() != helper.StringValue(remote.Label)
			descChanged := planImg.Description.ValueString() != helper.StringValue(remote.Description)

			if labelChanged || descChanged {
				opts := linodego.ImageShareGroupUpdateImageOptions{}
				if labelChanged {
					opts.Label = planImg.Label.ValueStringPointer()
				}
				if descChanged {
					opts.Description = planImg.Description.ValueStringPointer()
				}
				// Use SharedID for the update call
				toUpdate = append(toUpdate, struct {
					id   string
					opts linodego.ImageShareGroupUpdateImageOptions
				}{
					id:   remote.SharedID,
					opts: opts,
				})
			}
		}
	}

	// Detect removals
	for _, remote := range remoteMap {
		if _, exists := planMap[remote.PrivateID]; !exists {
			toRemove = append(toRemove, remote.SharedID)
		}
	}

	// Apply changes
	if len(toAdd) > 0 {
		opts := linodego.ImageShareGroupAddImagesOptions{Images: toAdd}
		if _, err := client.ImageShareGroupAddImages(ctx, imageShareGroupID, opts); err != nil {
			resp.Diagnostics.AddError("Failed to add images", err.Error())
			return
		}
	}

	for _, u := range toUpdate {
		if _, err := client.ImageShareGroupUpdateImageShareEntry(ctx, imageShareGroupID, u.id, u.opts); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to update image %s", u.id), err.Error())
			return
		}
	}

	for _, id := range toRemove {
		if err := client.ImageShareGroupRemoveImage(ctx, imageShareGroupID, id); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to remove image %s", id), err.Error())
			return
		}
	}

	// Update Image Share Group
	labelChanged := !plan.Label.Equal(state.Label)
	descChanged := !plan.Description.Equal(state.Description)
	if labelChanged || descChanged {
		updateOpts := linodego.ImageShareGroupUpdateOptions{
			Label:       plan.Label.ValueStringPointer(),
			Description: plan.Description.ValueStringPointer(),
		}
		if _, err := client.UpdateImageShareGroup(ctx, imageShareGroupID, updateOpts); err != nil {
			resp.Diagnostics.AddError("Failed to update share group", err.Error())
			return
		}
	}

	// Refresh and persist final state
	sg, err = client.GetImageShareGroup(ctx, imageShareGroupID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to re-fetch share group", err.Error())
		return
	}

	finalImages, err := client.ImageShareGroupListImageShareEntries(ctx, imageShareGroupID, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to re-fetch share group images", err.Error())
		return
	}

	// Flatten and update state
	state.FlattenImageShareGroup(sg, false)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert images into framework list
	var imageModels []ImageShareAttributesModel
	for _, img := range finalImages {
		sourceID := ""
		if img.ImageSharing.SharedBy != nil && img.ImageSharing.SharedBy.SourceImageID != nil {
			sourceID = *img.ImageSharing.SharedBy.SourceImageID
		}

		imageModels = append(imageModels, ImageShareAttributesModel{
			ID:          types.StringValue(sourceID),
			Label:       stringToPointerValue(img.Label),
			Description: stringToPointerValue(img.Description),
		})
	}

	listVal, diag := types.ListValueFrom(ctx, imageShareGroupImage.Type(), imageModels)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Images = listVal

	// Normalize null/unknown -> []
	if state.Images.IsNull() || state.Images.IsUnknown() {
		state.Images = types.ListValueMust(imageShareGroupImage.Type(), []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	imageShareGroupID := helper.FrameworkSafeInt64ToInt(state.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Meta.Client

	err := client.DeleteImageShareGroup(ctx, imageShareGroupID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Delete Image Share Group", err.Error())
		return
	}
}

// This is a custom struct to store a subset of im_ImageShare row data for use in detecting/running updates and deletions
// for the Images in an Image Share Group
type imageData struct {
	PrivateID   string  // original private ID from API
	SharedID    string  // source_image_id from ImageSharing.SharedBy.SourceImageID
	Label       *string // optional label
	Description *string // optional description
}

func stringToPointerValue(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
