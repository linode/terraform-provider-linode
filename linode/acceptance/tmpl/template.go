package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

// ProviderNoPoll is used to configure the provider to disable instance
// polling.
func ProviderNoPoll(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"provider_no_poll", nil,
	)
}
