package acceptance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type FrameworkProviderWithClient struct {
	linode.FrameworkProvider
	client *linodego.Client
}

func NewFrameworkProviderWithClient(
	client *linodego.Client,
) provider.Provider {
	return &FrameworkProviderWithClient{
		FrameworkProvider: *TestAccFrameworkProvider,
		client:            client,
	}
}

func (fp *FrameworkProviderWithClient) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	// Call parent configure func
	fp.FrameworkProvider.Configure(ctx, req, resp)

	resp.ResourceData.(*helper.FrameworkProviderMeta).Client = fp.client
	resp.DataSourceData.(*helper.FrameworkProviderMeta).Client = fp.client
}
