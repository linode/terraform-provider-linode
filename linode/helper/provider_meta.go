package helper

import (
	"fmt"

	"github.com/linode/linodego"
)

func GetFwClientWithUserAgent(
	resourceOrDataSourceUAComment string,
	meta *FrameworkProviderMeta,
) *linodego.Client {
	client := meta.Client
	client.SetUserAgent(generateUserAgent(meta.ProviderUserAgent, resourceOrDataSourceUAComment))

	return &client
}

func GetSDKClientWithUserAgent(
	resourceOrDataSourceUAComment string,
	meta *ProviderMeta,
) linodego.Client {
	client := meta.Client
	client.SetUserAgent(generateUserAgent(meta.ProviderUserAgent, resourceOrDataSourceUAComment))

	return client
}

func generateUserAgent(providerUserAgent, resourceOrDataSourceUAComment string) string {
	if resourceOrDataSourceUAComment == "" {
		return providerUserAgent
	}

	return providerUserAgent + fmt.Sprintf(" (%s)", resourceOrDataSourceUAComment)
}
