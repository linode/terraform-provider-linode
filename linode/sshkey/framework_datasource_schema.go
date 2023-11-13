package sshkey

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var SSHKeyAttributes = map[string]schema.Attribute{
	"label": schema.StringAttribute{
		Description: "The label of the Linode SSH Key.",
		Required:    true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(0, 64),
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
	"id": schema.StringAttribute{
		Description: "A unique identifier for this datasource.",
		Optional:    true,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: SSHKeyAttributes,
}
