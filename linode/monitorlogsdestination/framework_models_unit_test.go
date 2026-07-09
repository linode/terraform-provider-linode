//go:build unit

package monitorlogsdestination

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestFlattenLogsDestination_AkamaiObjectStorage(t *testing.T) {
	createdTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)

	dest := &linodego.LogsDestination{
		ID:        42,
		Label:     "my-destination",
		Type:      linodego.LogsDestinationTypeAkamaiObjectStorage,
		Status:    "active",
		CreatedBy: "user1",
		UpdatedBy: "user2",
		Created:   &createdTime,
		Updated:   &updatedTime,
		Version:   3,
		Details: linodego.LogsDestinationDetails{
			AccessKeyID: "AKID123",
			BucketName:  "my-bucket",
			Host:        "us-east-1.linodeobjects.com",
			Path:        "/logs",
		},
	}

	m := &LogsDestinationResourceModel{
		AkamaiObjectStorageDetails: &LogsDestinationAkamaiDetailsModel{
			// Simulate preserved write-only secret from state
			AccessKeySecret: types.StringValue("super-secret"),
		},
	}

	diags := m.FlattenLogsDestination(context.Background(), dest, false)
	assert.Empty(t, diags)

	assert.Equal(t, types.StringValue("42"), m.ID)
	assert.Equal(t, types.StringValue("my-destination"), m.Label)
	assert.Equal(t, types.StringValue("akamai_object_storage"), m.Type)
	assert.Equal(t, types.StringValue("active"), m.Status)
	assert.Equal(t, types.StringValue("user1"), m.CreatedBy)
	assert.Equal(t, types.StringValue("user2"), m.UpdatedBy)
	assert.Equal(t, types.Int64Value(3), m.Version)

	assert.NotNil(t, m.AkamaiObjectStorageDetails)
	assert.Equal(t, types.StringValue("AKID123"), m.AkamaiObjectStorageDetails.AccessKeyID)
	assert.Equal(t, types.StringValue("my-bucket"), m.AkamaiObjectStorageDetails.BucketName)
	assert.Equal(t, types.StringValue("us-east-1.linodeobjects.com"), m.AkamaiObjectStorageDetails.Host)
	assert.Equal(t, types.StringValue("/logs"), m.AkamaiObjectStorageDetails.Path)

	// Write-only field must NOT be overwritten by FlattenLogsDestination
	assert.Equal(t, types.StringValue("super-secret"), m.AkamaiObjectStorageDetails.AccessKeySecret)

	assert.Nil(t, m.CustomHTTPSDetails)
}

func TestFlattenLogsDestination_CustomHTTPS(t *testing.T) {
	createdTime := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, time.March, 2, 0, 0, 0, 0, time.UTC)

	dest := &linodego.LogsDestination{
		ID:        99,
		Label:     "https-dest",
		Type:      linodego.LogsDestinationTypeCustomHTTPS,
		Status:    "active",
		CreatedBy: "admin",
		UpdatedBy: "admin",
		Created:   &createdTime,
		Updated:   &updatedTime,
		Version:   1,
		Details: linodego.LogsDestinationDetails{
			EndpointURL:     "https://logs.example.com/ingest",
			ContentType:     "application/json",
			DataCompression: "gzip",
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: "basic",
			},
			ClientCertificateDetails: &linodego.LogsDestinationClientCertificateDetails{
				TLSHostname: "logs.example.com",
			},
			CustomHeaders: []linodego.LogsDestinationCustomHTTPSHeader{
				{Name: "X-Custom-Header", Value: ""},
			},
		},
	}

	m := &LogsDestinationResourceModel{
		CustomHTTPSDetails: &LogsDestinationCustomHTTPSDetailsModel{
			Authentication: &LogsDestinationAuthModel{
				Username: types.StringValue("myuser"),
				Password: types.StringValue("mypass"),
			},
			ClientCertificateDetails: &LogsDestinationClientCertDetailsModel{
				ClientCACertificate: types.StringValue("ca-cert"),
				ClientCertificate:   types.StringValue("client-cert"),
				ClientPrivateKey:    types.StringValue("private-key"),
			},
			CustomHeaders: []LogsDestinationCustomHeaderModel{
				{
					Name:  types.StringValue("X-Custom-Header"),
					Value: types.StringValue("secret-header-value"),
				},
			},
		},
	}

	diags := m.FlattenLogsDestination(context.Background(), dest, false)
	assert.Empty(t, diags)

	assert.NotNil(t, m.CustomHTTPSDetails)
	d := m.CustomHTTPSDetails

	assert.Equal(t, types.StringValue("https://logs.example.com/ingest"), d.EndpointURL)
	assert.Equal(t, types.StringValue("application/json"), d.ContentType)
	assert.Equal(t, types.StringValue("gzip"), d.DataCompression)

	// Auth type is returned by API and should be updated
	assert.NotNil(t, d.Authentication)
	assert.Equal(t, types.StringValue("basic"), d.Authentication.Type)

	// Write-only fields must be preserved from state
	assert.Equal(t, types.StringValue("myuser"), d.Authentication.Username)
	assert.Equal(t, types.StringValue("mypass"), d.Authentication.Password)

	// TLS hostname is returned by API
	assert.NotNil(t, d.ClientCertificateDetails)
	assert.Equal(t, types.StringValue("logs.example.com"), d.ClientCertificateDetails.TLSHostname)

	// Write-only cert fields must be preserved from state
	assert.Equal(t, types.StringValue("ca-cert"), d.ClientCertificateDetails.ClientCACertificate)
	assert.Equal(t, types.StringValue("client-cert"), d.ClientCertificateDetails.ClientCertificate)
	assert.Equal(t, types.StringValue("private-key"), d.ClientCertificateDetails.ClientPrivateKey)

	// Custom header name is returned, sensitive value preserved from state
	assert.Len(t, d.CustomHeaders, 1)
	assert.Equal(t, types.StringValue("X-Custom-Header"), d.CustomHeaders[0].Name)
	assert.Equal(t, types.StringValue("secret-header-value"), d.CustomHeaders[0].Value)

	assert.Nil(t, m.AkamaiObjectStorageDetails)
}

func TestParseLogsDestination_DataSource(t *testing.T) {
	createdTime := time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, time.June, 2, 0, 0, 0, 0, time.UTC)

	dest := &linodego.LogsDestination{
		ID:        7,
		Label:     "ds-dest",
		Type:      linodego.LogsDestinationTypeAkamaiObjectStorage,
		Status:    "active",
		CreatedBy: "ds-user",
		UpdatedBy: "ds-user",
		Created:   &createdTime,
		Updated:   &updatedTime,
		Version:   2,
		Details: linodego.LogsDestinationDetails{
			AccessKeyID: "DSKEY",
			BucketName:  "ds-bucket",
			Host:        "eu-central-1.linodeobjects.com",
			Path:        "/ds-logs",
		},
	}

	m := &LogsDestinationDataSourceModel{}
	m.ParseLogsDestination(dest)

	assert.Equal(t, types.Int64Value(7), m.ID)
	assert.Equal(t, types.StringValue("ds-dest"), m.Label)
	assert.Equal(t, types.StringValue("akamai_object_storage"), m.Type)
	assert.Equal(t, types.StringValue("active"), m.Status)
	assert.Equal(t, types.StringValue("ds-user"), m.CreatedBy)
	assert.Equal(t, types.StringValue("ds-user"), m.UpdatedBy)
	assert.Equal(t, types.Int64Value(2), m.Version)

	assert.NotNil(t, m.Details)
	assert.Equal(t, types.StringValue("DSKEY"), m.Details.AccessKeyID)
	assert.Equal(t, types.StringValue("ds-bucket"), m.Details.BucketName)
	assert.Equal(t, types.StringValue("eu-central-1.linodeobjects.com"), m.Details.Host)
	assert.Equal(t, types.StringValue("/ds-logs"), m.Details.Path)
}

func TestParseLogsDestination_DataSource_CustomHTTPS(t *testing.T) {
	createdTime := time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)
	updatedTime := time.Date(2024, time.August, 2, 0, 0, 0, 0, time.UTC)

	dest := &linodego.LogsDestination{
		ID:        55,
		Label:     "https-ds-dest",
		Type:      linodego.LogsDestinationTypeCustomHTTPS,
		Status:    "active",
		CreatedBy: "creator",
		UpdatedBy: "updater",
		Created:   &createdTime,
		Updated:   &updatedTime,
		Version:   3,
		Details: linodego.LogsDestinationDetails{
			EndpointURL:     "https://logs.example.com/ingest",
			ContentType:     "application/json",
			DataCompression: "gzip",
			Authentication: &linodego.LogsDestinationCustomHTTPSAuthDetails{
				Type: "basic",
			},
			ClientCertificateDetails: &linodego.LogsDestinationClientCertificateDetails{
				TLSHostname: "logs.example.com",
			},
		},
	}

	m := &LogsDestinationDataSourceModel{}
	m.ParseLogsDestination(dest)

	assert.Equal(t, types.Int64Value(55), m.ID)
	assert.Equal(t, types.StringValue("https-ds-dest"), m.Label)
	assert.Equal(t, types.StringValue("custom_https"), m.Type)
	assert.Equal(t, types.StringValue("active"), m.Status)
	assert.Equal(t, types.StringValue("creator"), m.CreatedBy)
	assert.Equal(t, types.StringValue("updater"), m.UpdatedBy)
	assert.Equal(t, types.Int64Value(3), m.Version)

	assert.NotNil(t, m.Details)
	assert.Equal(t, types.StringValue("https://logs.example.com/ingest"), m.Details.EndpointURL)
	assert.Equal(t, types.StringValue("application/json"), m.Details.ContentType)
	assert.Equal(t, types.StringValue("gzip"), m.Details.DataCompression)
	assert.Equal(t, types.StringValue("basic"), m.Details.AuthenticationType)
	assert.Equal(t, types.StringValue("logs.example.com"), m.Details.TLSHostname)

	// OBJ storage fields must be empty for a custom_https destination.
	assert.True(t, m.Details.AccessKeyID.IsNull())
	assert.True(t, m.Details.BucketName.IsNull())
	assert.True(t, m.Details.Host.IsNull())
}
