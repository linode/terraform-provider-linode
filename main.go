package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/linode/terraform-provider-linode/v3/linode"
	"github.com/linode/terraform-provider-linode/v3/version"
)

func main() {
	ctx := context.Background()

	// Disable the baked-in timestamp in favor of tflog
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSDKProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		linode.Provider().GRPCProvider,
	)
	if err != nil {
		log.Fatal("failed to upgrade SDKv2 GRPC provider:", err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(
			linode.CreateFrameworkProvider(version.ProviderVersion),
		),
		func() tfprotov6.ProviderServer {
			return upgradedSDKProvider
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/linode/linode",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
