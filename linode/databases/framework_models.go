package databases

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
)

// DatabaseModel represents a single Database object.
type DatabaseModel struct {
	ID              types.Int64    `tfsdk:"id"`
	AllowList       []types.String `tfsdk:"allow_list"`
	ClusterSize     types.Int64    `tfsdk:"cluster_size"`
	Created         types.String   `tfsdk:"created"`
	Encrypted       types.Bool     `tfsdk:"encrypted"`
	Engine          types.String   `tfsdk:"engine"`
	HostPrimary     types.String   `tfsdk:"host_primary"`
	HostSecondary   types.String   `tfsdk:"host_secondary"`
	InstanceURI     types.String   `tfsdk:"instance_uri"`
	Label           types.String   `tfsdk:"label"`
	Region          types.String   `tfsdk:"region"`
	ReplicationType types.String   `tfsdk:"replication_type"`
	SSLConnection   types.Bool     `tfsdk:"ssl_connection"`
	Status          types.String   `tfsdk:"status"`
	Type            types.String   `tfsdk:"type"`
	Updated         types.String   `tfsdk:"updated"`
	Version         types.String   `tfsdk:"version"`
}

// DatabaseFilterModel describes the Terraform resource data model to match the
// resource schema.
type DatabaseFilterModel struct {
	ID        types.String                     `tfsdk:"id"`
	Filters   frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order     types.String                     `tfsdk:"order"`
	OrderBy   types.String                     `tfsdk:"order_by"`
	Databases []DatabaseModel                  `tfsdk:"databases"`
}

// parseDatabases parses the given list of regions into the `databases` model attribute.
func (model *DatabaseFilterModel) parseDatabases(databases []linodego.Database) {
	parseDB := func(db linodego.Database) DatabaseModel {
		var m DatabaseModel
		m.ID = types.Int64Value(int64(db.ID))
		m.AllowList = helper.StringSliceToFramework(db.AllowList)
		m.ClusterSize = types.Int64Value(int64(db.ClusterSize))
		m.Encrypted = types.BoolValue(db.Encrypted)
		m.Engine = types.StringValue(db.Engine)
		m.HostPrimary = types.StringValue(db.Hosts.Primary)
		m.HostSecondary = types.StringValue(db.Hosts.Secondary)
		m.InstanceURI = types.StringValue(db.InstanceURI)
		m.Label = types.StringValue(db.Label)
		m.Region = types.StringValue(db.Region)
		m.ReplicationType = types.StringValue(db.ReplicationType)
		m.SSLConnection = types.BoolValue(db.SSLConnection)
		m.Status = types.StringValue(string(db.Status))
		m.Type = types.StringValue(db.Type)
		m.Version = types.StringValue(db.Version)

		if db.Created != nil {
			m.Created = types.StringValue(db.Created.Format(time.RFC3339))
		} else {
			m.Created = types.StringNull()
		}

		if db.Updated != nil {
			m.Updated = types.StringValue(db.Updated.Format(time.RFC3339))
		} else {
			m.Updated = types.StringNull()
		}

		return m
	}

	result := make([]DatabaseModel, len(databases))

	for i, db := range databases {
		result[i] = parseDB(db)
	}

	model.Databases = result
}
