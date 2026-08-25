package sshkey

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

const (
	SSHKeyLabelRegex        = "^[a-zA-Z0-9_\\-\\s]*$"
	SSHKeyLabelErrorMessage = "Labels may only contain letters, numbers, dashes, underscores, and spaces."
)

// frameworkDatasourceSchema supports lookup by either label or id.
// Both attributes are Optional+Computed so existing label-based configs keep
// working, id-based lookup is available, and the unset selector is filled from
// the API response.
var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label of the Linode SSH Key to select. " +
				"Exactly one of `label` or `id` must be specified.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("label"),
					path.MatchRoot("id"),
				),
				stringvalidator.LengthBetween(1, 64),
				helper.RegexMatches(SSHKeyLabelRegex, SSHKeyLabelErrorMessage),
			},
		},
		"id": schema.StringAttribute{
			Description: "The ID of the SSH Key to select. " +
				"Exactly one of `label` or `id` must be specified.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(
					path.MatchRoot("label"),
					path.MatchRoot("id"),
				),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"ssh_key": schema.StringAttribute{
			Description: "The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			CustomType:  timetypes.RFC3339Type{},
			Description: "The date this key was added.",
			Computed:    true,
		},
	},
}
