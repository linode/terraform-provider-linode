package account

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"email": {
		Type: schema.TypeString,
		Description: "The email address for this Account, for account management communications, " +
			"and may be used for other communications as configured.",
		Computed: true,
	},
	"first_name": {
		Type:        schema.TypeString,
		Description: "The first name of the person associated with this Account.",
		Computed:    true,
	},
	"last_name": {
		Type:        schema.TypeString,
		Description: "The last name of the person associated with this Account.",
		Computed:    true,
	},
	"company": {
		Type:        schema.TypeString,
		Description: "The company name associated with this Account.",
		Computed:    true,
	},
	"address_1": {
		Type:        schema.TypeString,
		Description: "First line of this Account's billing address.",
		Computed:    true,
	},
	"address_2": {
		Type:        schema.TypeString,
		Description: "Second line of this Account's billing address.",
		Computed:    true,
	},
	"phone": {
		Type:        schema.TypeString,
		Description: "The phone number associated with this Account.",
		Computed:    true,
	},
	"city": {
		Type:        schema.TypeString,
		Description: "The city for this Account's billing address.",
		Computed:    true,
	},
	"state": {
		Type: schema.TypeString,
		Description: "If billing address is in the United States, " +
			"this is the State portion of the Account's billing address. " +
			"If the address is outside the US, this is the Province associated with the Account's billing address.",
		Computed: true,
	},
	"country": {
		Type:        schema.TypeString,
		Description: "The two-letter country code of this Account's billing address.",
		Computed:    true,
	},
	"zip": {
		Type:        schema.TypeString,
		Description: "The zip code of this Account's billing address.",
		Computed:    true,
	},
	"balance": {
		Type:        schema.TypeInt,
		Description: "This Account's balance, in US dollars.",
		Computed:    true,
	},
}
