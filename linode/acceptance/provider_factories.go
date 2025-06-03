package acceptance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"linode": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSDKProvider, err := tf5to6server.UpgradeServer(
			context.Background(),
			TestAccSDKv2Providers["linode"].GRPCProvider,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to upgrade SDKv2 GRPC provider: %w", err)
		}

		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(
				TestAccFrameworkProvider,
			),
			func() tfprotov6.ProviderServer {
				return upgradedSDKProvider
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}

var HttpExternalProviders = map[string]resource.ExternalProvider{
	"http": {
		Source: "hashicorp/http",
	},
}
