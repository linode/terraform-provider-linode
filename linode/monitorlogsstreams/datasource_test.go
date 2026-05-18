//go:build integration || monitorlogsstreams

package monitorlogsstreams_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsstreams/tmpl"
)

func TestAccDataSourceMonitorLogsStreams_basic(t *testing.T) {
	acceptance.LongRunningTest(t)

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.linode_monitor_logs_streams.foobar",
						tfjsonpath.New("streams"), knownvalue.NotNull()),
				},
			},
		},
	})
}
