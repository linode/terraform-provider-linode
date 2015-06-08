package main

import (
	"github.com/btobolaskiterraform-linode"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: linode.Provider,
	})
}
