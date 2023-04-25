package acceptance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"linode": func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()
		providers := []func() tfprotov5.ProviderServer{
			TestAccProviders["linode"].GRPCProvider,
			providerserver.NewProtocol5(
				TestAccFrameworkProvider,
			),
		}

		muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
