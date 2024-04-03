//go:build integration || accountlogin

package accountlogin_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/accountlogin/tmpl"
)

func TestAccDataSourceLinodeAccountLogin_basic(t *testing.T) {
	acceptance.OptInTest(t)
	t.Parallel()

	resourceName := "data.linode_account_login.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	logins, err := client.ListLogins(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to list logins: %s", err)
	}

	login := logins[0]
	accountID := login.ID

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing login.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", strconv.Itoa(login.ID)),
					resource.TestCheckResourceAttr(resourceName, "ip", login.IP),
					resource.TestCheckResourceAttr(resourceName, "username", login.Username),
					resource.TestCheckResourceAttr(resourceName, "datetime", login.Datetime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "restricted", strconv.FormatBool(login.Restricted)),
					resource.TestCheckResourceAttr(resourceName, "status", login.Status),
				),
			},
		},
	})
}
