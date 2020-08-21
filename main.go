package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/linode/terraform-provider-linode/linode"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: linode.Provider,
	})
}
