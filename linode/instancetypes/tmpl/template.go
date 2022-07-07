package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_basic", nil)
}

func DataSubstring(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_substring", nil)
}

func DataRegex(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"instance_types_data_regex", nil)
}
