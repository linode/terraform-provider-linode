package obj

import (
	"context"
	"strings"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

const (
	REGION_CLUSTER_REQUIRE_REPLACEMENT_FUNC_DESCRIPTION = "Require replacement if region or " +
		"cluster has been changed and the change is not a migration from a cluster to " +
		"an equivalent region or from a region to an equivalent cluster"
)

func requireReplacementIfClusterOrRegionSemanticallyChanged(
	ctx context.Context,
	sr planmodifier.StringRequest,
	rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse,
) {
	var regionPlan, clusterPlan, regionState, clusterState types.String
	sr.Plan.GetAttribute(ctx, path.Root("cluster"), &clusterPlan)
	sr.Plan.GetAttribute(ctx, path.Root("region"), &regionPlan)
	sr.State.GetAttribute(ctx, path.Root("cluster"), &clusterState)
	sr.State.GetAttribute(ctx, path.Root("region"), &regionState)

	if !regionState.IsNull() && regionPlan.IsNull() && !clusterPlan.IsNull() &&
		strings.HasPrefix(clusterPlan.ValueString(), regionState.ValueString()) {
		// the region changed to an equivalent cluster
		return
	}

	if !clusterState.IsNull() && clusterPlan.IsNull() && !regionPlan.IsNull() &&
		strings.HasPrefix(clusterState.ValueString(), regionPlan.ValueString()) {
		// the cluster changed to an equivalent region
		return
	}

	rrifr.RequiresReplace = true
}

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the object.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"bucket": schema.StringAttribute{
			Description: "The target bucket to put this object in.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"cluster": schema.StringAttribute{
			Description: "The target cluster that the bucket is in.",
			DeprecationMessage: "The cluster attribute has been deprecated, please consider switching to the region attribute. " +
				"For example, a cluster value of `us-mia-1` can be translated to a region value of `us-mia`.",
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					requireReplacementIfClusterOrRegionSemanticallyChanged,
					REGION_CLUSTER_REQUIRE_REPLACEMENT_FUNC_DESCRIPTION,
					REGION_CLUSTER_REQUIRE_REPLACEMENT_FUNC_DESCRIPTION,
				),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("region")),
			},
		},
		"region": schema.StringAttribute{
			Description: "The target region that the bucket is in.",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					requireReplacementIfClusterOrRegionSemanticallyChanged,
					REGION_CLUSTER_REQUIRE_REPLACEMENT_FUNC_DESCRIPTION,
					REGION_CLUSTER_REQUIRE_REPLACEMENT_FUNC_DESCRIPTION,
				),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("cluster")),
			},
		},
		"key": schema.StringAttribute{
			Description: "The name of the uploaded object.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"secret_key": schema.StringAttribute{
			Description: "The REQUIRED S3 secret key with access to the target bucket. " +
				"If not specified with the resource, you must provide its value by configuring the obj_secret_key, " +
				"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
			Optional:  true,
			Sensitive: true,
		},
		"access_key": schema.StringAttribute{
			Description: "The REQUIRED S3 access key with access to the target bucket. " +
				"If not specified with the resource, you must provide its value by configuring the obj_access_key, " +
				"or, opting-in generating it implicitly at apply-time using obj_use_temp_keys at provider-level.",
			Optional: true,
		},
		"content": schema.StringAttribute{
			Description: "The contents of the Object to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content_base64"),
					path.MatchRoot("source"),
				),
			},
		},
		"content_base64": schema.StringAttribute{
			Description: "The base64 contents of the Object to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content"),
					path.MatchRoot("source"),
				),
			},
		},
		"source": schema.StringAttribute{
			Description: "The source file to upload.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("content"),
					path.MatchRoot("content_base64"),
				),
			},
		},
		"acl": schema.StringAttribute{
			Description: "The ACL config given to this object.",
			Optional:    true,
			Computed:    true,
			Default: stringdefault.StaticString(
				string(s3types.ObjectCannedACLPrivate),
			),
			Validators: []validator.String{
				stringvalidator.OneOf(
					helper.StringAliasSliceToStringSlice(
						s3types.ObjectCannedACLPrivate.Values(),
					)...,
				),
			},
		},
		"cache_control": schema.StringAttribute{
			Description: "This cache_control configuration of this object.",
			Optional:    true,
		},
		"content_disposition": schema.StringAttribute{
			Description: "The content disposition configuration of this object.",
			Optional:    true,
		},
		"content_encoding": schema.StringAttribute{
			Description: "The encoding of the content of this object.",
			Optional:    true,
		},
		"content_language": schema.StringAttribute{
			Description: "The language metadata of this object.",
			Optional:    true,
		},
		"content_type": schema.StringAttribute{
			Description: "The MIME type of the content.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"endpoint": schema.StringAttribute{
			Description: "The endpoint for the bucket used for s3 connections.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"etag": schema.StringAttribute{
			Description: "The specific version of this object.",
			Optional:    true,
			Computed:    true,
		},
		"force_destroy": schema.BoolAttribute{
			Description: "Whether the object should bypass deletion restrictions.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"metadata": schema.MapAttribute{
			Description: "The metadata of this object",
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
			Default:     helper.EmptyMapDefault(types.StringType),
		},
		"version_id": schema.StringAttribute{
			Description: "The version ID of this object.",
			Computed:    true,
		},
		"website_redirect": schema.StringAttribute{
			Description: "The website redirect location of this object.",
			Optional:    true,
		},
	},
}
