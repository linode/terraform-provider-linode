package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_data_basic", nil)
}
