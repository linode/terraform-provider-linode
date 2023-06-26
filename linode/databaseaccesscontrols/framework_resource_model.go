package databaseaccesscontrols

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceModel describes the Terraform resource model to match the
// resource schema.
type ResourceModel struct {
	DatabaseID   types.Int64    `tfsdk:"database_id"`
	DatabaseType types.String   `tfsdk:"database_type"`
	AllowList    []types.String `tfsdk:"allow_list"`

	ID types.String `tfsdk:"id"`
}
