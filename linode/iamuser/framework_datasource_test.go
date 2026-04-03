//go:build integration || iamuser

package iamuser_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/iamuser/tmpl"
)

func TestAccDataSourceIAMUser_basic(t *testing.T) {
	t.Parallel()

	// IAM Tests need to be opted into, iam accounts do not support all existing user endpoints as they will be replacing some of them
	acceptance.OptInTest(t)

	resName := "data.linode_iam_user.test_iam_user"
	username := acctest.RandomWithPrefix("tf_test")
	email := username + "@example.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, username, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "account_access.0", "account_event_viewer"),
				),
			},
		},
	})
}
