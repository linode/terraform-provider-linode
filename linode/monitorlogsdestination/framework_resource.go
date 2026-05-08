package monitorlogsdestination

import (
	"context"
	"fmt"
	"strconv"

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
				Name:   "linode_monitor_logs_destination",
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

	var data LogsDestinationResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.LogsDestinationCreateOptions{
		Label: data.Label.ValueString(),
		Type:  linodego.LogsDestinationType(data.Type.ValueString()),
	}

	switch linodego.LogsDestinationType(data.Type.ValueString()) {
	case linodego.LogsDestinationTypeAkamaiObjectStorage:
		if data.AkamaiObjectStorageDetails == nil {
			resp.Diagnostics.AddError(
				"Missing akamai_object_storage_details",
				"akamai_object_storage_details block is required when type is akamai_object_storage.",
			)
			return
		}
		d := data.AkamaiObjectStorageDetails
		details := linodego.LogsDestinationDetailsCreateOptions{
			AccessKeyID:     d.AccessKeyID.ValueString(),
			AccessKeySecret: d.AccessKeySecret.ValueString(),
			BucketName:      d.BucketName.ValueString(),
			Host:            d.Host.ValueString(),
		}
		if !d.Path.IsNull() && !d.Path.IsUnknown() {
			path := d.Path.ValueString()
			details.Path = &path
		}
		createOpts.Details = details

	case linodego.LogsDestinationTypeCustomHTTPS:
		if data.CustomHTTPSDetails == nil {
			resp.Diagnostics.AddError(
				"Missing custom_https_details",
				"custom_https_details block is required when type is custom_https.",
			)
			return
		}
		d := data.CustomHTTPSDetails
		details := linodego.LogsDestinationCustomHTTPSDetailsCreateOptions{
			EndpointURL:     d.EndpointURL.ValueString(),
			ContentType:     d.ContentType.ValueString(),
			DataCompression: d.DataCompression.ValueString(),
		}
		if d.Authentication != nil {
			auth := &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: linodego.LogsDestinationCustomHTTPSAuthType(d.Authentication.Type.ValueString()),
			}
			if !d.Authentication.Username.IsNull() && !d.Authentication.Username.IsUnknown() {
				auth.Details = &linodego.LogsDestinationCustomHTTPSBasicAuthDetails{
					Username: d.Authentication.Username.ValueString(),
					Password: d.Authentication.Password.ValueString(),
				}
			}
			details.Authentication = auth
		}
		if d.ClientCertificateDetails != nil {
			details.ClientCertificateDetails = &linodego.LogsDestinationClientCertificateDetails{
				TLSHostname:         d.ClientCertificateDetails.TLSHostname.ValueString(),
				ClientCACertificate: d.ClientCertificateDetails.ClientCACertificate.ValueString(),
				ClientCertificate:   d.ClientCertificateDetails.ClientCertificate.ValueString(),
				ClientPrivateKey:    d.ClientCertificateDetails.ClientPrivateKey.ValueString(),
			}
		}
		if len(d.CustomHeaders) > 0 {
			headers := make([]linodego.LogsDestinationCustomHTTPSHeader, len(d.CustomHeaders))
			for i, h := range d.CustomHeaders {
				headers[i] = linodego.LogsDestinationCustomHTTPSHeader{
					Name:  h.Name.ValueString(),
					Value: h.Value.ValueString(),
				}
			}
			details.CustomHeaders = headers
		}
		createOpts.Details = details
	}

	tflog.Debug(ctx, "client.CreateLogsDestination(...)")
	dest, err := client.CreateLogsDestination(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create logs destination.", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FlattenLogsDestination(ctx, dest, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(strconv.Itoa(dest.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	tflog.Debug(ctx, "Read "+r.Config.Name)

	var data LogsDestinationResourceModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if helper.FrameworkAttemptRemoveResourceForEmptyID(ctx, data.ID, resp) {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", id)

	tflog.Trace(ctx, "client.GetLogsDestination(...)")
	dest, err := client.GetLogsDestination(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Logs destination no longer exists.",
				fmt.Sprintf(
					"Removing logs destination with ID %v from state because it no longer exists.",
					id,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to refresh the logs destination.", err.Error())
		return
	}

	resp.Diagnostics.Append(data.FlattenLogsDestination(ctx, dest, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	tflog.Debug(ctx, "Update "+r.Config.Name)

	client := r.Meta.Client
	var plan, state LogsDestinationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(state.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", id)

	updateOpts := linodego.LogsDestinationUpdateOptions{}

	if !state.Label.Equal(plan.Label) {
		updateOpts.Label = plan.Label.ValueString()
	}

	destType := linodego.LogsDestinationType(plan.Type.ValueString())
	switch destType {
	case linodego.LogsDestinationTypeAkamaiObjectStorage:
		if plan.AkamaiObjectStorageDetails != nil {
			d := plan.AkamaiObjectStorageDetails
			details := linodego.LogsDestinationDetailsUpdateOptions{
				AccessKeyID:     d.AccessKeyID.ValueString(),
				AccessKeySecret: d.AccessKeySecret.ValueString(),
				BucketName:      d.BucketName.ValueString(),
				Host:            d.Host.ValueString(),
			}
			if !d.Path.IsNull() && !d.Path.IsUnknown() {
				path := d.Path.ValueString()
				details.Path = &path
			}
			updateOpts.Details = details
		}

	case linodego.LogsDestinationTypeCustomHTTPS:
		if plan.CustomHTTPSDetails != nil {
			d := plan.CustomHTTPSDetails
			details := linodego.LogsDestinationCustomHTTPSDetailsUpdateOptions{
				EndpointURL:     d.EndpointURL.ValueString(),
				ContentType:     d.ContentType.ValueString(),
				DataCompression: d.DataCompression.ValueString(),
			}
			if d.Authentication != nil {
				auth := &linodego.LogsDestinationCustomHTTPSAuthDetails{
					Type: linodego.LogsDestinationCustomHTTPSAuthType(d.Authentication.Type.ValueString()),
				}
				if !d.Authentication.Username.IsNull() && !d.Authentication.Username.IsUnknown() {
					auth.Details = &linodego.LogsDestinationCustomHTTPSBasicAuthDetails{
						Username: d.Authentication.Username.ValueString(),
						Password: d.Authentication.Password.ValueString(),
					}
				}
				details.Authentication = auth
			}
			if d.ClientCertificateDetails != nil {
				details.ClientCertificateDetails = &linodego.LogsDestinationClientCertificateDetails{
					TLSHostname:         d.ClientCertificateDetails.TLSHostname.ValueString(),
					ClientCACertificate: d.ClientCertificateDetails.ClientCACertificate.ValueString(),
					ClientCertificate:   d.ClientCertificateDetails.ClientCertificate.ValueString(),
					ClientPrivateKey:    d.ClientCertificateDetails.ClientPrivateKey.ValueString(),
				}
			}
			if len(d.CustomHeaders) > 0 {
				headers := make([]linodego.LogsDestinationCustomHTTPSHeader, len(d.CustomHeaders))
				for i, h := range d.CustomHeaders {
					headers[i] = linodego.LogsDestinationCustomHTTPSHeader{
						Name:  h.Name.ValueString(),
						Value: h.Value.ValueString(),
					}
				}
				details.CustomHeaders = headers
			}
			updateOpts.Details = details
		}
	}

	tflog.Debug(ctx, "client.UpdateLogsDestination(...)")
	dest, err := client.UpdateLogsDestination(ctx, id, updateOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to update logs destination (%d).", id),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.FlattenLogsDestination(ctx, dest, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.CopyFrom(state, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	tflog.Debug(ctx, "Delete "+r.Config.Name)

	client := r.Meta.Client
	var data LogsDestinationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "client.DeleteLogsDestination(...)")
	err := client.DeleteLogsDestination(ctx, id)
	if err != nil {
		if linodego.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete logs destination (%d).", id),
			err.Error(),
		)
	}
}
