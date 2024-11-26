package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

func DataList(t *testing.T) string {
	return acceptance.ExecuteTemplate(t, "reserved_ips_data", nil)
}
