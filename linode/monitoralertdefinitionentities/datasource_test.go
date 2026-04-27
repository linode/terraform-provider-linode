//go:build integration || monitoralertdefinitionentities

package monitoralertdefinitionentities_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitoralertdefinitionentities/tmpl"
)

func TestAccDataSourceAlertDefinitionEntities_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_monitor_alert_definition_entities.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("service_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("alert_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("entities"), knownvalue.NotNull()),
				},
			},
		},
	})
}
