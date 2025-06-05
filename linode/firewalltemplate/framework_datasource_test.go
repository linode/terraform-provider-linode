//go:build integration || firewalltemplate

package firewalltemplate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalltemplate/tmpl"
)

const testTemplateDataName = "data.linode_firewall_template.test"

func TestAccDataSourceFirewalls_basic(t *testing.T) {
	t.Parallel()
	testSlug := "public"

	acceptance.RunTestWithRetries(t, 3, func(t *acceptance.WrappedT) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, testSlug),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("id"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("slug"), knownvalue.StringExact(testSlug)),
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("inbound"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("inbound_policy"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("outbound"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue(testTemplateDataName, tfjsonpath.New("outbound_policy"), knownvalue.NotNull()),
					},
				},
			},
		})
	})
}
