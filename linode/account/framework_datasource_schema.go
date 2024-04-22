package account

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"email": schema.StringAttribute{
			Description: "The email address for this Account, for account management communications, " +
				"and may be used for other communications as configured.",
			Computed: true,
		},
		"first_name": schema.StringAttribute{
			Description: "The first name of the person associated with this Account.",
			Computed:    true,
		},
		"last_name": schema.StringAttribute{
			Description: "The last name of the person associated with this Account.",
			Computed:    true,
		},
		"company": schema.StringAttribute{
			Description: "The company name associated with this Account.",
			Computed:    true,
		},
		"address_1": schema.StringAttribute{
			Description: "First line of this Account's billing address.",
			Computed:    true,
		},
		"address_2": schema.StringAttribute{
			Description: "Second line of this Account's billing address.",
			Computed:    true,
		},
		"phone": schema.StringAttribute{
			Description: "The phone number associated with this Account.",
			Computed:    true,
		},
		"city": schema.StringAttribute{
			Description: "The city for this Account's billing address.",
			Computed:    true,
		},
		"state": schema.StringAttribute{
			Description: "If billing address is in the United States, " +
				"this is the State portion of the Account's billing address. " +
				"If the address is outside the US, this is the Province associated with the Account's billing address.",
			Computed: true,
		},
		"country": schema.StringAttribute{
			Description: "The two-letter country code of this Account's billing address.",
			Computed:    true,
		},
		"zip": schema.StringAttribute{
			Description: "The zip code of this Account's billing address.",
			Computed:    true,
		},
		"balance": schema.Float64Attribute{
			Description: "This Account's balance, in US dollars.",
			Computed:    true,
		},
		"id": schema.StringAttribute{
			Description: "The Email of the Account.",
			Computed:    true,
		},
	},
}
