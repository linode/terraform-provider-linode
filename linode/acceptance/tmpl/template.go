package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

// ProviderNoPoll is used to configure the provider to disable instance
// polling.
func ProviderNoPoll(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"provider_no_poll", nil,
	)
}
