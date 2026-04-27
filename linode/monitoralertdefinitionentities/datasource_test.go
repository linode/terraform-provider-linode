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

//var channelID int
//
//func init() {
//	client, err := acceptance.GetTestClient()
//	if err != nil {
//		log.Fatal(fmt.Errorf("Error getting client: %s", err))
//	}
//
//	channels, err := client.ListAlertChannels(context.Background(), nil)
//	if err != nil {
//		log.Fatal(fmt.Errorf("error listing alert channels: %s", err))
//	}
//	if len(channels) < 1 {
//		log.Fatal(fmt.Errorf("at least one alert channel is required"))
//	}
//
//	channelID = channels[0].ID
//}

func TestAccDataSourceAlertDefinitionEntities_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_monitor_alert_definition_entities.foobar"
	// alertLabel := acctest.RandomWithPrefix("tf-test")
	// alertChannels := channelID

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
