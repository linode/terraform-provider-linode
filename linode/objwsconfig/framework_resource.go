package objwsconfig

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_object_storage_website_config",
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
	tflog.Debug(ctx, "Create linode_object_storage_website_config")
	var data ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  data.Bucket.ValueString(),
		"cluster": data.Cluster.ValueString(),
	})

	tflog.Info(ctx, "Creating the bucket website config")

	websiteDomain, err := GetS3WebsiteDomain(ctx, r.Meta.Client, data.Cluster.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read website domain for cluster", err.Error())
		return
	}

	err = putBucketWebsite(ctx, r.Meta.Client, data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create bucket website config", err.Error())
		return
	}

	data.ComputeWebsiteEndpoint(websiteDomain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read linode_object_storage_website_config")

	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  data.Bucket.ValueString(),
		"cluster": data.Cluster.ValueString(),
	})

	s3client, err := s3ConnectionFromData(ctx, r.Meta.Client, data)
	if err != nil {
		resp.Diagnostics.AddError("Unable to refresh the bucket website config", err.Error())
		return
	}

	websiteDomain, err := GetS3WebsiteDomain(ctx, r.Meta.Client, data.Cluster.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read website domain for cluster", err.Error())
		return
	}

	tflog.Debug(ctx, "Fetching the bucket website config")

	output, err := s3client.GetBucketWebsite(ctx, &s3.GetBucketWebsiteInput{
		Bucket: aws.String(data.Bucket.ValueString()),
	})

	if err == nil {
		data.FlattenBucketWebsite(output, false)
		data.ComputeWebsiteEndpoint(websiteDomain)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else if isNotFoundError(err) {
		resp.Diagnostics.AddWarning(
			"The bucket website config does not exist",
			fmt.Sprintf(
				"Removing bucket website config with bucket %v from state because it no longer exists",
				data.Bucket.ValueString(),
			),
		)
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.AddError("Unable to refresh the bucket website config", err.Error())
	}
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update linode_object_storage_website_config")
	var plan, state ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  state.Bucket.ValueString(),
		"cluster": state.Cluster.ValueString(),
	})

	shouldUpdate := !state.IndexDocument.Equal(plan.IndexDocument) || !state.ErrorDocument.Equal(plan.ErrorDocument)
	if shouldUpdate {
		tflog.Info(ctx, "Updating bucket website config")
		err := putBucketWebsite(ctx, r.Meta.Client, plan)
		if err != nil {
			resp.Diagnostics.AddError("Failed to update bucket website config", err.Error())
			return
		}
	}

	plan.CopyFrom(state, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete linode_object_storage_website_config")
	var data ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"bucket":  data.Bucket.ValueString(),
		"cluster": data.Cluster.ValueString(),
	})

	tflog.Info(ctx, "Deleting the bucket website config")

	err := deleteBucketWebsite(ctx, r.Meta.Client, data)
	if err != nil && !isNotFoundError(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete the bucket website config for bucket %s", data.Bucket.ValueString()),
			err.Error(),
		)
	}
}
