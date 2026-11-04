//go:build integration || natgateway

package natgateway_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgateway/tmpl"
)

const (
	natAutoTrue     = "linode_natgateway.explicit_autoscaling_true"
	natAutoFalse    = "linode_natgateway.explicit_autoscaling_false"
	natDefaultEmpty = "linode_natgateway.default_empty_addresses"
	natDefaultPop   = "linode_natgateway.default_populated_addresses"
	natComplex      = "linode_natgateway.complex"
)

// TestAccResourceNATGateway_update creates the four NAT Gateway variants, changes their
// labels, and swaps the manual gateway's address. Each update must happen in place.
func TestAccResourceNATGateway_update(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The firewall label is derived from this value and must begin with a letter.
	label := "tf-" + acctest.RandString(10)
	updatedLabel := "tf-" + acctest.RandString(10)

	// natGatewayStateChecks returns the expected state for each variant. base is the
	// label prefix used in the template for this step.
	natGatewayStateChecks := func(base string) []statecheck.StateCheck {
		return []statecheck.StateCheck{
			// Explicit autoscaling: the configured value is kept in state.
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("use_autoscaling"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-autoscaling-true")),
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(0)),

			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("use_autoscaling"), knownvalue.Bool(false)),
			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-autoscaling-false")),
			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),

			// Implied autoscaling: not configured, so not stored (the API does not return it).
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("use_autoscaling"), knownvalue.Null()),
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-default-empty")),
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(0)),

			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("use_autoscaling"), knownvalue.Null()),
			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-default-populated")),
			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create all four variants.
				Config:            tmpl.Basic(t, label, region),
				ConfigStateChecks: natGatewayStateChecks(label),
			},
			{
				// Update labels only. Every variant must update in place (no replacement),
				// and use_autoscaling / addresses must be unchanged.
				Config: tmpl.Basic(t, updatedLabel, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(natAutoTrue, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natAutoFalse, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natDefaultEmpty, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natDefaultPop, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: natGatewayStateChecks(updatedLabel),
			},
			{
				// Swap the manual gateway's address from address 1 to address 3. Address 3 is
				// not used by any other gateway. Must be an in-place update.
				Config: tmpl.AddressSwap(t, updatedLabel, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(natAutoFalse, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: natGatewayStateChecks(updatedLabel),
			},
			{
				// Assign an address to the autoscaling-true gateway. The API must reject this.
				// Address 1 is free here because the swap moved the manual gateway to address 3.
				Config:      tmpl.AutoscaleAddressUpdate(t, updatedLabel, region),
				ExpectError: regexp.MustCompile(`Addresses cannot be modified when autoscaling is enabled`),
			},
		},
	})
}

// TestAccResourceNATGateway_complex creates a manual (autoscaling false) NAT Gateway with one
// reserved address, attached to a VPC subnet created in the same config.
func TestAccResourceNATGateway_complex(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityNATGateway}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The label must begin with a letter.
	label := "tf-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("label"), knownvalue.StringExact(label+"-complex")),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("use_autoscaling"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("region"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("default_ports_per_interface"), knownvalue.Int64Exact(4096)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("address_autoscale_max"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("portset_assignments"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("portset_capacity"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("type"), knownvalue.StringExact("vpc_subnet")),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("url"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("vpc_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("vpc_label"), knownvalue.StringExact(label)),
				},
			},
			{
				// Changing vpc_subnet_id forces a replacement. The gateway is still attached
				// to its original subnet when it is destroyed, so Delete must detach it.
				Config: tmpl.ComplexAltSubnet(t, label, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(natComplex, plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet").AtMapKey("label"), knownvalue.StringExact(label+"-alt")),
				},
			},
		},
	})
}
