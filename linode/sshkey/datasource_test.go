//go:build integration || sshkey

package sshkey_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/sshkey/tmpl"
)

func TestAccDataSourceSSHKey_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_sshkey.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, acceptance.PublicKeyMaterial),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("ssh_key"), knownvalue.StringExact(acceptance.PublicKeyMaterial)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
			{
				Config: tmpl.DataByID(t, label, acceptance.PublicKeyMaterial),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("ssh_key"), knownvalue.StringExact(acceptance.PublicKeyMaterial)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("created"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDataSourceSSHKey_notFound(t *testing.T) {
	t.Parallel()

	missingLabel := acctest.RandomWithPrefix("tf_test") + "-missing"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.DataMissing(t, missingLabel),
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(missingLabel) + " was not found"),
			},
			{
				// Use a high ID that should not exist for this account.
				Config:      tmpl.DataMissingByID(t, "999999999"),
				ExpectError: regexp.MustCompile(`(?i)(not found|failed to get ssh key)`),
			},
		},
	})
}
