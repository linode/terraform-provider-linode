package networkingipassignment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkingIPModel struct {
	ID          types.String      `tfsdk:"id"`
	Region      types.String      `tfsdk:"region"`
	Assignments []AssignmentModel `tfsdk:"assignments"`
}

type AssignmentModel struct {
	Address  types.String `tfsdk:"address"`
	LinodeID types.Int64  `tfsdk:"linode_id"`
}
