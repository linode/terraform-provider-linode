package databases

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseDatabases(t *testing.T) {
	mockDB1 := linodego.Database{
		ID:              123,
		Status:          linodego.DatabaseStatusActive,
		Label:           "example-db-1",
		Region:          "us-east",
		Type:            "g6-standard-1",
		Engine:          "mysql",
		Version:         "8.0",
		Encrypted:       false,
		AllowList:       []string{"203.0.113.1/32", "192.0.1.0/24"},
		Hosts:           linodego.DatabaseHost{Primary: "primary.example.com", Secondary: "secondary.example.com"},
		InstanceURI:     "mysql://user:pass@primary.example.com:3306/db",
		ReplicationType: "standard",
		SSLConnection:   true,
		Created:         &time.Time{},
		Updated:         &time.Time{},
	}

	mockDB2 := linodego.Database{
		ID:              456,
		Status:          linodego.DatabaseStatusProvisioning,
		Label:           "example-db-2",
		Region:          "us-central",
		Type:            "g6-standard-2",
		Engine:          "postgresql",
		Version:         "13",
		Encrypted:       true,
		AllowList:       []string{"10.0.0.1/32"},
		Hosts:           linodego.DatabaseHost{Primary: "primary-pg.example.com", Secondary: "secondary-pg.example.com"},
		InstanceURI:     "postgresql://user:pass@primary-pg.example.com:5432/db",
		ReplicationType: "standard",
		SSLConnection:   false,
		Created:         &time.Time{},
		Updated:         &time.Time{},
	}

	mockDatabases := []linodego.Database{mockDB1, mockDB2}

	model := &DatabaseFilterModel{}
	model.parseDatabases(mockDatabases)

	assert.Len(t, model.Databases, 2)

	// Database 1 Assertions
	assert.Equal(t, types.Int64Value(123), model.Databases[0].ID)
	assert.Equal(t, types.StringValue("active"), model.Databases[0].Status)
	assert.Equal(t, types.StringValue("example-db-1"), model.Databases[0].Label)

	// Database 2 Assertions
	assert.Equal(t, types.Int64Value(456), model.Databases[1].ID)
	assert.Equal(t, types.StringValue("provisioning"), model.Databases[1].Status)
	assert.Equal(t, types.StringValue("example-db-2"), model.Databases[1].Label)
}
