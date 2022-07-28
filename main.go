package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/linode/terraform-provider-linode/linode"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Debug:        true,
		ProviderFunc: linode.Provider,
	})
}
