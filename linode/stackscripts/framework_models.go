package stackscripts

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscript"
)

// StackscriptFilterModel describes the Terraform resource data model to match the
// resource schema.
type StackscriptFilterModel struct {
	ID           types.String                     `tfsdk:"id"`
	Latest       types.Bool                       `tfsdk:"latest"`
	Filters      frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order        types.String                     `tfsdk:"order"`
	OrderBy      types.String                     `tfsdk:"order_by"`
	Stackscripts []stackscript.StackScriptModel   `tfsdk:"stackscripts"`
}

func (data *StackscriptFilterModel) parseStackscripts(
	stackscripts []linodego.Stackscript,
) diag.Diagnostics {
	result := make([]stackscript.StackScriptModel, len(stackscripts))

	for i := range stackscripts {
		var stackscript stackscript.StackScriptModel
		diags := stackscript.FlattenStackScript(&stackscripts[i], false)
		if diags != nil {
			return diags
		}
		result[i] = stackscript
	}

	data.Stackscripts = result
	return nil
}
