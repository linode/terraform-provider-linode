package placementgroupassignment

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type idFormat struct {
	PGID     int `json:"pg_id"`
	LinodeID int `json:"linode_ids"`
}

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
	idData := idFormat{
		PGID:     pg.ID,
		LinodeID: expectedLinodeID,
	}

	m.ID = helper.KeepOrUpdateString(m.ID, buildID(idData, diags), preserveKnown)
	m.PlacementGroupID = helper.KeepOrUpdateInt64(m.PlacementGroupID, int64(pg.ID), preserveKnown)

	if pgHasID(pg, expectedLinodeID) {
		m.LinodeID = helper.KeepOrUpdateInt64(m.LinodeID, int64(expectedLinodeID), preserveKnown)
	} else {
		m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, types.Int64Null(), preserveKnown)
	}
}

func (m *PGAssignmentModel) CopyFrom(
	other PGAssignmentModel,
	preserveKnown bool,
) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.PlacementGroupID = helper.KeepOrUpdateValue(m.PlacementGroupID, other.PlacementGroupID, preserveKnown)
	m.LinodeID = helper.KeepOrUpdateValue(m.LinodeID, other.LinodeID, preserveKnown)
}

func buildID(idData idFormat, diags *diag.Diagnostics) string {
	renderedJSON, err := json.Marshal(idData)
	if err != nil {
		diags.AddError(
			"Failed to marshal JSON for ID",
			err.Error(),
		)
		return ""
	}

	return base64.StdEncoding.EncodeToString(renderedJSON)
}

func parseID(idStr string, diags *diag.Diagnostics) idFormat {
	var idData idFormat

	idStrDecoded, err := base64.StdEncoding.DecodeString(idStr)
	if err != nil {
		diags.AddError(
			"Failed to decode ID base64",
			err.Error(),
		)
		return idData
	}

	if err := json.Unmarshal(idStrDecoded, &idData); err != nil {
		diags.AddError(
			"Failed to marshal JSON for ID",
			err.Error(),
		)
	}

	return idData
}

func pgHasID(pg linodego.PlacementGroup, id int) bool {
	for _, member := range pg.Members {
		if member.LinodeID == id {
			return true
		}
	}

	return false
}
