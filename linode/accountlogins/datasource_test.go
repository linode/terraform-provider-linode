package accountlogins_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/accountlogins/tmpl"
)

func TestAccDataSourceAccountLogins_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_logins.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.ip"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.restricted"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.datetime"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "logins.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceAccountLogins_filterByRestricted(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_logins.foobar"

	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	logins, err := client.ListLogins(context.TODO(), nil)
	randIndex := rand.Intn(len(logins))
	login := logins[randIndex]

	username := login.Username
	ip := login.IP
	restricted := login.Restricted

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing login.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterRestricted(t, username, ip, restricted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.ip"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.username"),
					resource.TestCheckResourceAttr(resourceName, "logins.0.restricted", strconv.FormatBool(restricted)),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.datetime"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "logins.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceAccountLogins_filterByUsername(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_logins.foobar"

	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	logins, err := client.ListLogins(context.TODO(), nil)
	randIndex := rand.Intn(len(logins))
	login := logins[randIndex]

	username := login.Username
	ip := login.IP
	restricted := login.Restricted

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing login.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterUsername(t, username, ip, restricted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.ip"),
					resource.TestCheckResourceAttr(resourceName, "logins.0.username", username),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.restricted"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.datetime"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "logins.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceAccountLogins_filterByIP(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_account_logins.foobar"

	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	logins, err := client.ListLogins(context.TODO(), nil)
	randIndex := rand.Intn(len(logins))
	login := logins[randIndex]

	username := login.Username
	ip := login.IP
	restricted := login.Restricted

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing login.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterIP(t, username, ip, restricted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.id"),
					resource.TestCheckResourceAttr(resourceName, "logins.0.ip", ip),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.restricted"),
					resource.TestCheckResourceAttrSet(resourceName, "logins.0.datetime"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "logins.#", 0),
				),
			},
		},
	})
}
