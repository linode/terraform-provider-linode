package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_config_data_basic", nil)
}
