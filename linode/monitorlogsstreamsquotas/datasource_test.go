//go:build integration || monitorlogsstreamsquotas

package monitorlogsstreamsquotas_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/monitorlogsstreamsquotas/tmpl"
)

func TestAccDataSourceMonitorLogsStreamQuotas_basic(t *testing.T) {
	t.Parallel()

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatalf("failed to get test client: %s", err)
	}

	quotas, err := client.ListLogStreamQuotas(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list log stream quotas: %s", err)
	}

	if len(quotas) == 0 {
		t.Skip("No available Monitor Logs Stream quotas for testing.")
	}

	const dataName = "data.linode_monitor_logs_stream_quotas.all"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas"),
						knownvalue.ListSizeExact(len(quotas)),
					),
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_limit"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataName,
						tfjsonpath.New("quotas").AtSliceIndex(0).AtMapKey("quota_type"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
