package obj

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	READ_PERMISSION       = "read_only"
	READ_WRITE_PERMISSION = "read_write"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_object_storage_object",
				IDType: types.StringType,
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
	config := r.Meta.Config

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, plan)

	plan.ComputeEndpointIfUnknown(ctx, client, &resp.Diagnostics)

	s3client, teardownKeys := getS3ClientFromModel(
		ctx, client, config, plan, READ_WRITE_PERMISSION, &resp.Diagnostics,
	)

	if teardownKeys != nil {
		defer teardownKeys()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	fwPutObject(ctx, plan, s3client, &resp.Diagnostics)

	// Add resource to TF states earlier to prevent
	// dangling resources (resources created but not managed by TF)
	AddObjectResource(ctx, resp, plan)

	RefreshObject(ctx, &plan, s3client, &resp.Diagnostics, nil, true)

	// IDs should always be overridden during creation (see #1085)
	// TODO: Remove when Crossplane empty string ID issue is resolved
	plan.GenerateObjectStorageObjectID(true, false)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func RefreshObject(
	ctx context.Context,
	data *ResourceModel,
	s3client *s3.Client,
	diags *diag.Diagnostics,
	removeResource func(context.Context),
	preserveKnown bool,
) {
	tflog.Debug(ctx, "enter RefreshObject")

	bucket := data.Bucket.ValueString()
	key := data.Key.ValueString()

	if diags.HasError() {
		return
	}

	headObjectInput := &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	tflog.Debug(ctx, "getting object header", map[string]any{"HeadObjectInput": headObjectInput})
	headOutput, err := s3client.HeadObject(
		ctx,
		headObjectInput,
	)
	if err != nil {
		if helper.IsObjNotFoundErr(err) && removeResource != nil {
			removeResource(ctx)
			diags.AddWarning(
				"Object Not Found",
				"couldn't find the bucket or object, removing the object from the TF state",
			)
		}
		diags.AddError("Failed to Refresh the Object", err.Error())
	}

	data.FlattenObject(*headOutput, preserveKnown)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	client := r.Meta.Client
	config := r.Meta.Config

	var state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = populateLogAttributes(ctx, state)

	// TODO: cleanup when Crossplane fixes it
	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, state.ID, resp) {
		return
	}

	s3client, teardownKeys := getS3ClientFromModel(
		ctx, client, config, state, READ_PERMISSION, &resp.Diagnostics,
	)

	if teardownKeys != nil {
		defer teardownKeys()
	}

	if resp.Diagnostics.HasError() {
		if newDiags := deleteBucketNotFound(resp.Diagnostics); len(newDiags) < len(resp.Diagnostics) {
			resp.Diagnostics = newDiags

			resp.Diagnostics.AddWarning(
				"The Object No Longer Exists",
				fmt.Sprintf(
					"Removing Object Storage Object %q from state because it no longer exists",
					state.ID.ValueString(),
				),
			)

			resp.State.RemoveResource(ctx)
		}
		return
	}

	RefreshObject(ctx, &state, s3client, &resp.Diagnostics, resp.State.RemoveResource, false)

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

	ctx = populateLogAttributes(ctx, state)
	client := r.Meta.Client
	config := r.Meta.Config

	plan.ComputeEndpointIfUnknown(ctx, client, &resp.Diagnostics)

	s3client, teardownKeys := getS3ClientFromModel(
		ctx, client, config, plan, READ_WRITE_PERMISSION, &resp.Diagnostics,
	)

	if teardownKeys != nil {
		defer teardownKeys()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if (!plan.ETag.IsUnknown() && !plan.ETag.Equal(state.ETag)) ||
		!plan.CacheControl.Equal(state.CacheControl) ||
		!plan.ContentBase64.Equal(state.ContentBase64) ||
		!plan.ContentDisposition.Equal(state.ContentDisposition) ||
		!plan.ContentEncoding.Equal(state.ContentEncoding) ||
		!plan.ContentLanguage.Equal(state.ContentLanguage) ||
		!plan.ContentType.Equal(state.ContentType) ||
		!plan.Content.Equal(state.Content) ||
		!plan.Metadata.Equal(state.Metadata) ||
		!plan.Source.Equal(state.Source) ||
		!plan.WebsiteRedirect.Equal(state.WebsiteRedirect) {

		fwPutObject(ctx, plan, s3client, &resp.Diagnostics)
	}

	RefreshObject(ctx, &plan, s3client, &resp.Diagnostics, nil, true)

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

	ctx = populateLogAttributes(ctx, state)

	client := r.Meta.Client
	config := r.Meta.Config

	s3client, teardownKeys := getS3ClientFromModel(
		ctx, client, config, state, READ_WRITE_PERMISSION, &resp.Diagnostics,
	)

	if teardownKeys != nil {
		defer teardownKeys()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	force := state.ForceDestroy.ValueBool()
	bucket := state.Bucket.ValueString()
	key := state.Key.ValueString()

	if !state.VersionID.IsNull() {
		tflog.Debug(ctx, "versioning was enabled for this object, deleting all versions and delete markers")

		err := helper.DeleteAllObjectVersionsAndDeleteMarkers(ctx, s3client, bucket, key, force, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Delete All Object Versions and Deletion Markers in the Versioned Bucket",
				err.Error(),
			)
		}
	} else {
		tflog.Debug(ctx, "versioning was disabled for this object, simply delete the object")

		err := deleteObject(ctx, s3client, bucket, strings.TrimPrefix(key, "/"), "", force)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Delete the Object", err.Error())
		}
	}
}

func populateLogAttributes(ctx context.Context, model ResourceModel) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":            model.Bucket.ValueString(),
		"region_or_cluster": model.GetRegionOrCluster(ctx),
		"object_key":        model.Key.ValueString(),
	})
}
