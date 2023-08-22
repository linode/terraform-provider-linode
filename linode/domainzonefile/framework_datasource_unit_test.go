//go:build unit

package domainzonefile

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseDomainZoneFile(t *testing.T) {
	mockZoneFileData := linodego.DomainZoneFile{
		ZoneFile: []string{
			"; example.com [123]",
			"$TTL 864000",
			"@  IN  SOA  ns1.linode.com. user.example.com. 2021000066 14400 14400 1209600 86400",
			"@    NS  ns1.linode.com.",
			"@    NS  ns2.linode.com.",
			"@    NS  ns3.linode.com.",
			"@    NS  ns4.linode.com.",
			"@    NS  ns5.linode.com.",
		},
	}

	var data DataSourceModel
	diags := data.parseDomainZoneFile(context.Background(), &mockZoneFileData)
	assert.False(t, diags.HasError(), "Error parsing domain zone file")

	assert.Empty(t, data.DomainID, "DomainID should be 0/nil")
	for _, line := range mockZoneFileData.ZoneFile {
		assert.Contains(t, data.ZoneFile.String(), line, "ZoneFile content doesn't contain expected line")
	}

	idJSON, _ := json.Marshal(&mockZoneFileData)
	assert.Equal(t, types.StringValue(string(idJSON)), data.ID, "ID doesn't match")
}
