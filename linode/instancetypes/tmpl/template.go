package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_basic", nil)
}

func DataSubstring(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_substring", nil)
}

func DataRegex(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_regex", nil)
}

func DataByClass(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_by_class", nil)
}
