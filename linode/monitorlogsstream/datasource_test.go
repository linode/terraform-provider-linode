//go:build integration || monitorlogsstream

package monitorlogsstream_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/monitorlogsstream/tmpl"
)

func TestAccDataSourceMonitorLogsStream_notFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.DataNotFound(t),
				ExpectError: regexp.MustCompile(`(?i)404|not found`),
			},
		},
	})
}
