package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"testing"
)

type TemplateData struct {
	Label string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_basic", TemplateData{
			Label: label,
		})
}

func RouteTarget(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_route_target", TemplateData{
			Label: label,
		})
}

func NoID(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"ipv6range_no_id", nil)
}
