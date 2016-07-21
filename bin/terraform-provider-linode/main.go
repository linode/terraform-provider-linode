package main

import (
	"github.com/btobolaski/terraform-provider-linode"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: linode.Provider,
	})
}
