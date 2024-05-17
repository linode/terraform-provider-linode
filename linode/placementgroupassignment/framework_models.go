package placementgroupassignment

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type PGAssignmentModel struct {
	ID               types.String `tfsdk:"id"`
	PlacementGroupID types.Int64  `tfsdk:"placement_group_id"`
	LinodeID         types.Int64  `tfsdk:"linode_id"`
	CompliantOnly    types.Bool   `tfsdk:"compliant_only"`
}

func (m *PGAssignmentModel) Flatten(
	pg linodego.PlacementGroup,
	expectedLinodeID int,
	preserveKnown bool,
	diags *diag.Diagnostics,
) {
	m.ID = helper.KeepOrUpdateString(m.ID, buildID(pg.ID, expectedLinodeID, diags), preserveKnown)
	m.PlacementGroupID = helper.KeepOrUpdateInt64(m.PlacementGroupID, int64(pg.ID), preserveKnown)
	m.LinodeID = helper.KeepOrUpdateInt64(m.LinodeID, int64(expectedLinodeID), preserveKnown)
}

func (m *PGAssignmentModel) CopyFrom(
	other PGAssignmentModel,
	preserveKnown bool,
) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.PlacementGroupID = helper.KeepOrUpdateValue(m.PlacementGroupID, other.PlacementGroupID, preserveKnown)
	m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, other.LinodeID, preserveKnown)
}

func (m *PGAssignmentModel) GetIDComponents(diags *diag.Diagnostics) (pgID int, linodeID int) {
	pgID = helper.FrameworkSafeInt64ToInt(m.PlacementGroupID.ValueInt64(), diags)
	linodeID = helper.FrameworkSafeInt64ToInt(m.LinodeID.ValueInt64(), diags)
	return
}

func buildID(pgID, linodeID int, diags *diag.Diagnostics) string {
	renderedJSON, err := json.Marshal(
		struct {
			PGID     int `json:"pg_id"`
			LinodeID int `json:"linode_id"`
		}{
			PGID:     pgID,
			LinodeID: linodeID,
		},
	)
	if err != nil {
		diags.AddError(
			"Failed to marshal JSON for ID",
			err.Error(),
		)
		return ""
	}

	return base64.StdEncoding.EncodeToString(renderedJSON)
}
