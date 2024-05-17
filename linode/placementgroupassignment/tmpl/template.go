package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	AssignmentExists bool
}

func Basic(t *testing.T, label, region string, assignmentExists bool) string {
	return acceptance.ExecuteTemplate(t,
		"placement_group_assignment_basic", TemplateData{
			Label:            label,
			Region:           region,
			AssignmentExists: assignmentExists,
		})
}
