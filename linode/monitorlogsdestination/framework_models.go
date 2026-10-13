package monitorlogsdestination

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// LogsDestinationResourceModel is the model for the resource.
type LogsDestinationResourceModel struct {
	ID        types.String      `tfsdk:"id"`
	Label     types.String      `tfsdk:"label"`
	Type      types.String      `tfsdk:"type"`
	Status    types.String      `tfsdk:"status"`
	CreatedBy types.String      `tfsdk:"created_by"`
	UpdatedBy types.String      `tfsdk:"updated_by"`
	Created   timetypes.RFC3339 `tfsdk:"created"`
	Updated   timetypes.RFC3339 `tfsdk:"updated"`
	Version   types.Int64       `tfsdk:"version"`

	AkamaiObjectStorageDetails *LogsDestinationAkamaiDetailsModel      `tfsdk:"akamai_object_storage_details"`
	CustomHTTPSDetails         *LogsDestinationCustomHTTPSDetailsModel `tfsdk:"custom_https_details"`
}

// LogsDestinationAkamaiDetailsModel holds details for an akamai_object_storage destination.
type LogsDestinationAkamaiDetailsModel struct {
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"` // write-only, not returned by API
	BucketName      types.String `tfsdk:"bucket_name"`
	Host            types.String `tfsdk:"host"`
	Path            types.String `tfsdk:"path"`
}

// LogsDestinationCustomHTTPSDetailsModel holds details for a custom_https destination.
type LogsDestinationCustomHTTPSDetailsModel struct {
	EndpointURL              types.String                           `tfsdk:"endpoint_url"`
	ContentType              types.String                           `tfsdk:"content_type"`
	DataCompression          types.String                           `tfsdk:"data_compression"`
	Authentication           *LogsDestinationAuthModel              `tfsdk:"authentication"`
	ClientCertificateDetails *LogsDestinationClientCertDetailsModel `tfsdk:"client_certificate_details"`
	CustomHeaders            []LogsDestinationCustomHeaderModel     `tfsdk:"custom_headers"`
}

// LogsDestinationAuthModel holds authentication config for a custom_https destination.
type LogsDestinationAuthModel struct {
	Type     types.String `tfsdk:"type"`
	Username types.String `tfsdk:"username"` // write-only, not returned by API
	Password types.String `tfsdk:"password"` // write-only, not returned by API
}

// LogsDestinationClientCertDetailsModel holds TLS client certificate details.
type LogsDestinationClientCertDetailsModel struct {
	TLSHostname         types.String `tfsdk:"tls_hostname"`
	ClientCACertificate types.String `tfsdk:"client_ca_certificate"` // write-only
	ClientCertificate   types.String `tfsdk:"client_certificate"`    // write-only
	ClientPrivateKey    types.String `tfsdk:"client_private_key"`    // write-only
}

// LogsDestinationCustomHeaderModel holds a single custom HTTP header.
type LogsDestinationCustomHeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"` // sensitive, may not be returned by API
}

// LogsDestinationDataSourceModel is the model for the single-item data source.
type LogsDestinationDataSourceModel struct {
	ID        types.Int64                      `tfsdk:"id"`
	Label     types.String                     `tfsdk:"label"`
	Type      types.String                     `tfsdk:"type"`
	Status    types.String                     `tfsdk:"status"`
	CreatedBy types.String                     `tfsdk:"created_by"`
	UpdatedBy types.String                     `tfsdk:"updated_by"`
	Created   timetypes.RFC3339                `tfsdk:"created"`
	Updated   timetypes.RFC3339                `tfsdk:"updated"`
	Version   types.Int64                      `tfsdk:"version"`
	Details   *LogsDestinationFlatDetailsModel `tfsdk:"details"`
}

// LogsDestinationFlatDetailsModel is a flattened view of details for the data source.
// Only non-sensitive fields returned by the API are included.
type LogsDestinationFlatDetailsModel struct {
	// akamai_object_storage fields
	AccessKeyID types.String `tfsdk:"access_key_id"`
	BucketName  types.String `tfsdk:"bucket_name"`
	Host        types.String `tfsdk:"host"`
	Path        types.String `tfsdk:"path"`
	// custom_https fields
	EndpointURL        types.String `tfsdk:"endpoint_url"`
	ContentType        types.String `tfsdk:"content_type"`
	DataCompression    types.String `tfsdk:"data_compression"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	TLSHostname        types.String `tfsdk:"tls_hostname"`
}

// FlattenLogsDestination populates the resource model from an API response.
// Write-only fields (secrets, keys, passwords) are intentionally NOT updated
// since the API never returns them — they are preserved from existing state.
func (m *LogsDestinationResourceModel) FlattenLogsDestination(
	_ context.Context,
	dest *linodego.LogsDestination,
	preserveKnown bool,
) (resultDiags diag.Diagnostics) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(dest.ID), preserveKnown)
	m.Label = helper.KeepOrUpdateString(m.Label, dest.Label, preserveKnown)
	m.Type = helper.KeepOrUpdateString(m.Type, string(dest.Type), preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(dest.Status), preserveKnown)
	m.CreatedBy = helper.KeepOrUpdateString(m.CreatedBy, dest.CreatedBy, preserveKnown)
	m.UpdatedBy = helper.KeepOrUpdateString(m.UpdatedBy, dest.UpdatedBy, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(dest.Created), preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(dest.Updated), preserveKnown)
	m.Version = helper.KeepOrUpdateInt64(m.Version, int64(dest.Version), preserveKnown)

	switch dest.Type {
	case linodego.LogsDestinationTypeAkamaiObjectStorage:
		if m.AkamaiObjectStorageDetails == nil {
			m.AkamaiObjectStorageDetails = &LogsDestinationAkamaiDetailsModel{}
		}
		d := m.AkamaiObjectStorageDetails
		d.AccessKeyID = helper.KeepOrUpdateString(d.AccessKeyID, dest.Details.AccessKeyID, preserveKnown)
		// AccessKeySecret is write-only — never updated from API response
		d.BucketName = helper.KeepOrUpdateString(d.BucketName, dest.Details.BucketName, preserveKnown)
		d.Host = helper.KeepOrUpdateString(d.Host, dest.Details.Host, preserveKnown)
		d.Path = helper.KeepOrUpdateString(d.Path, dest.Details.Path, preserveKnown)

	case linodego.LogsDestinationTypeCustomHTTPS:
		if m.CustomHTTPSDetails == nil {
			m.CustomHTTPSDetails = &LogsDestinationCustomHTTPSDetailsModel{}
		}
		d := m.CustomHTTPSDetails
		d.EndpointURL = helper.KeepOrUpdateString(d.EndpointURL, dest.Details.EndpointURL, preserveKnown)
		d.ContentType = helper.KeepOrUpdateString(d.ContentType, dest.Details.ContentType, preserveKnown)
		d.DataCompression = helper.KeepOrUpdateString(d.DataCompression, dest.Details.DataCompression, preserveKnown)

		if dest.Details.Authentication != nil {
			if d.Authentication == nil {
				d.Authentication = &LogsDestinationAuthModel{}
			}
			d.Authentication.Type = helper.KeepOrUpdateString(
				d.Authentication.Type,
				string(dest.Details.Authentication.Type),
				preserveKnown,
			)
			// Username and Password are write-only — never updated from API response
		}

		if dest.Details.ClientCertificateDetails != nil {
			if d.ClientCertificateDetails == nil {
				d.ClientCertificateDetails = &LogsDestinationClientCertDetailsModel{}
			}
			d.ClientCertificateDetails.TLSHostname = helper.KeepOrUpdateString(
				d.ClientCertificateDetails.TLSHostname,
				dest.Details.ClientCertificateDetails.TLSHostname,
				preserveKnown,
			)
			// client_ca_certificate, client_certificate, client_private_key are write-only — never updated
		}

		// Rebuild custom_headers — match by name to preserve sensitive values
		if len(dest.Details.CustomHeaders) > 0 {
			headers := make([]LogsDestinationCustomHeaderModel, len(dest.Details.CustomHeaders))
			for i, h := range dest.Details.CustomHeaders {
				headers[i] = LogsDestinationCustomHeaderModel{
					Name: types.StringValue(h.Name),
				}
				// Preserve the header value from state if API doesn't return it
				if h.Value != "" {
					headers[i].Value = types.StringValue(h.Value)
				} else if d.CustomHeaders != nil {
					for _, stateHeader := range d.CustomHeaders {
						if stateHeader.Name.ValueString() == h.Name {
							headers[i].Value = stateHeader.Value
							break
						}
					}
				}
			}
			d.CustomHeaders = headers
		} else if !preserveKnown {
			d.CustomHeaders = nil
		}
	}

	return resultDiags
}

// CopyFrom copies values from another model instance, respecting preserveKnown.
func (m *LogsDestinationResourceModel) CopyFrom(other LogsDestinationResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.Label = helper.KeepOrUpdateValue(m.Label, other.Label, preserveKnown)
	m.Type = helper.KeepOrUpdateValue(m.Type, other.Type, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.CreatedBy = helper.KeepOrUpdateValue(m.CreatedBy, other.CreatedBy, preserveKnown)
	m.UpdatedBy = helper.KeepOrUpdateValue(m.UpdatedBy, other.UpdatedBy, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.Version = helper.KeepOrUpdateValue(m.Version, other.Version, preserveKnown)

	if m.AkamaiObjectStorageDetails == nil && other.AkamaiObjectStorageDetails != nil {
		copied := *other.AkamaiObjectStorageDetails
		m.AkamaiObjectStorageDetails = &copied
	}
	if m.CustomHTTPSDetails == nil && other.CustomHTTPSDetails != nil {
		copied := *other.CustomHTTPSDetails
		m.CustomHTTPSDetails = &copied
	}
}

// ParseLogsDestination populates the data source model from an API response.
func (m *LogsDestinationDataSourceModel) ParseLogsDestination(dest *linodego.LogsDestination) {
	m.ID = types.Int64Value(int64(dest.ID))
	m.Label = types.StringValue(dest.Label)
	m.Type = types.StringValue(string(dest.Type))
	m.Status = types.StringValue(string(dest.Status))
	m.CreatedBy = types.StringValue(dest.CreatedBy)
	m.UpdatedBy = types.StringValue(dest.UpdatedBy)
	m.Created = timetypes.NewRFC3339TimePointerValue(dest.Created)
	m.Updated = timetypes.NewRFC3339TimePointerValue(dest.Updated)
	m.Version = types.Int64Value(int64(dest.Version))

	flat := &LogsDestinationFlatDetailsModel{}
	d := dest.Details

	stringOrNull := func(s string) types.String {
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	}

	flat.AccessKeyID = stringOrNull(d.AccessKeyID)
	flat.BucketName = stringOrNull(d.BucketName)
	flat.Host = stringOrNull(d.Host)
	flat.Path = stringOrNull(d.Path)
	flat.EndpointURL = stringOrNull(d.EndpointURL)
	flat.ContentType = stringOrNull(d.ContentType)
	flat.DataCompression = stringOrNull(d.DataCompression)

	// Handle nested structs gracefully
	flat.AuthenticationType = types.StringNull()
	if d.Authentication != nil && d.Authentication.Type != "" {
		flat.AuthenticationType = types.StringValue(string(d.Authentication.Type))
	}

	flat.TLSHostname = types.StringNull()
	if d.ClientCertificateDetails != nil && d.ClientCertificateDetails.TLSHostname != "" {
		flat.TLSHostname = types.StringValue(d.ClientCertificateDetails.TLSHostname)
	}

	m.Details = flat
}
