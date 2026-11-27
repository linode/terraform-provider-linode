//go:build unit

package objendpoints

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseObjectStorageEndpoint_IAFields(t *testing.T) {
	s3 := "us-iad-obj-ia0.linodeobjects.com"
	var m ObjectStorageEndpointModel
	m.parseObjectStorageEndpoint(linodego.ObjectStorageEndpoint{
		Region:              "us-iad",
		S3Endpoint:          &s3,
		EndpointType:        "E0",
		SupportedKeyClasses: []linodego.ObjectStorageKeyClass{"E", "IA"},
		Country:             "us",
		RegionLabel:         "Washington, DC",
	})

	assert.Equal(t, []types.String{types.StringValue("E"), types.StringValue("IA")}, m.SupportedKeyClasses)
	assert.Equal(t, types.StringValue("us"), m.Country)
	assert.Equal(t, types.StringValue("Washington, DC"), m.RegionLabel)
}

func TestParseObjectStorageEndpoint_IAFieldsAbsent(t *testing.T) {
	var m ObjectStorageEndpointModel
	m.parseObjectStorageEndpoint(linodego.ObjectStorageEndpoint{Region: "us-east"})

	assert.Nil(t, m.SupportedKeyClasses)
	assert.True(t, m.Country.IsNull())
	assert.True(t, m.RegionLabel.IsNull())
}
