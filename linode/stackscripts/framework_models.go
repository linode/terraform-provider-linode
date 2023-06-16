package stackscripts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/stackscript"
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
	ctx context.Context,
	stackscripts []linodego.Stackscript,
) diag.Diagnostics {
	result := make([]stackscript.StackScriptModel, len(stackscripts))

	for i := range stackscripts {
		var stackscript stackscript.StackScriptModel
		diags := stackscript.ParseComputedAttributes(ctx, &stackscripts[i])
		if diags != nil {
			return diags
		}

		diags = stackscript.ParseNonComputedAttributes(ctx, &stackscripts[i])
		if diags != nil {
			return diags
		}
		result[i] = stackscript
	}

	data.Stackscripts = result
	return nil
}
