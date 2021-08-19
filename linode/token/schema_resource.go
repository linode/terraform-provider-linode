package token

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Description: "The label of the Linode Token.",
		Optional:    true,
	},
	"scopes": {
		Type: schema.TypeString,
		Description: "The scopes this token was created with. These define what parts of the Account the " +
			"token can be used to access. Many command-line tools, such as the Linode CLI, require tokens with " +
			"access to *. Tokens with more restrictive scopes are generally more secure.",
		Required: true,
		ForceNew: true,
	},
	"expiry": {
		Type: schema.TypeString,
		Description: "When this token will expire. Personal Access Tokens cannot be renewed, so after " +
			"this time the token will be completely unusable and a new token will need to be generated. Tokens " +
			"may be created with 'null' as their expiry and will never expire unless revoked.",
		Optional:         true,
		ValidateFunc:     validDateTime,
		ForceNew:         true,
		DiffSuppressFunc: equivalentDate,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "The date and time this token was created.",
		Computed:    true,
	},
	"token": {
		Type:        schema.TypeString,
		Sensitive:   true,
		Description: "The token used to access the API.",
		Computed:    true,
	},
}

func equivalentDate(k, old, new string, d *schema.ResourceData) bool {
	if dtOld, err := time.Parse("2006-01-02T15:04:05", old); err != nil {
		log.Printf("[WARN] could not parse date %s: %s", old, err)
		return false
	} else if dtNew, err := time.Parse("2006-01-02T15:04:05", new); err != nil {
		log.Printf("[WARN] could not parse date %s: %s", new, err)
		return false
	} else {
		return dtOld.Equal(dtNew)
	}
}

func validDateTime(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if _, err := time.Parse("2006-01-02T15:04:05Z", v); err != nil {
		es = append(es, fmt.Errorf("expected %s to be a datetime, got %s", k, v))
	}

	return
}
