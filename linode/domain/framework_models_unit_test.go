package domain

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseDomain(t *testing.T) {
	mockDomain := &linodego.Domain{
		ID:          1234,
		Domain:      "example.org",
		Type:        "master",
		Status:      "active",
		Description: "description",
		SOAEmail:    "admin@example.org",
		RetrySec:    300,
		MasterIPs:   nil,
		AXfrIPs:     nil,
		Tags:        []string{"example tag", "another example"},
		ExpireSec:   300,
		RefreshSec:  300,
		TTLSec:      300,
	}

	domainModel := &DomainModel{}

	domainModel.parseDomain(mockDomain)

	assert.Equal(t, types.Int64Value(1234), domainModel.ID)
	assert.Equal(t, types.StringValue("example.org"), domainModel.Domain)
	assert.Equal(t, types.StringValue("master"), domainModel.Type)
	assert.Equal(t, types.StringValue("active"), domainModel.Status)
	assert.Equal(t, types.StringValue("description"), domainModel.Description)
	assert.Equal(t, types.StringValue("admin@example.org"), domainModel.SOAEmail)
	assert.Equal(t, types.Int64Value(300), domainModel.RetrySec)
	assert.Empty(t, domainModel.MasterIPs)
	assert.Empty(t, domainModel.AXFRIPs)
	assert.Contains(t, domainModel.Tags, types.StringValue("example tag"))
	assert.Contains(t, domainModel.Tags, types.StringValue("another example"))
	assert.Equal(t, types.Int64Value(300), domainModel.ExpireSec)
	assert.Equal(t, types.Int64Value(300), domainModel.RefreshSec)
	assert.Equal(t, types.Int64Value(300), domainModel.TTLSec)
}
