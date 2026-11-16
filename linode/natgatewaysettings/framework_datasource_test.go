//go:build integration || natgatewaysettings

package natgatewaysettings_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgatewaysettings/tmpl"
)

func TestAccDataSourceNATGatewaySettings_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_natgateway_settings.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("allowed_ports_per_interface"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("maximum_autoscaling_addresses_per_natgateway"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("maximum_reserved_addresses_per_natgateway"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}
