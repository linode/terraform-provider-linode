package linode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/version"
	"golang.org/x/oauth2"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

// DefaultLinodeURL is the Linode APIv4 URL to use.
const DefaultLinodeURL = "https://api.linode.com/v4"

// Config represents the Linode provider configuration.
type Config struct {
	AccessToken string
	APIURL      string
	APIVersion  string
	UAPrefix    string

	terraformVersion string

	SkipInstanceReadyPoll        bool
	MinRetryDelayMilliseconds    int
	MaxRetryDelayMilliseconds    int
	EventPollMilliseconds        int
	LKEEventPollMilliseconds     int
	LKENodeReadyPollMilliseconds int
}

// Client returns a fully initialized Linode client.
func (c *Config) Client() linodego.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.AccessToken})
	oauthTransport := &oauth2.Transport{
		Source: tokenSource,
	}
	loggingTransport := logging.NewTransport("Linode", oauthTransport)

	oauth2Client := &http.Client{
		Transport: loggingTransport,
	}
	client := linodego.NewClient(oauth2Client)

	tfUserAgent := terraformUserAgent(c.terraformVersion)
	userAgent := strings.TrimSpace(fmt.Sprintf("%s terraform-provider-linode/%s",
		tfUserAgent, version.ProviderVersion))
	if c.UAPrefix != "" {
		userAgent = c.UAPrefix + " " + userAgent
	}
	client.SetUserAgent(userAgent)

	if c.APIURL != "" {
		client.SetBaseURL(c.APIURL)
	} else if len(c.APIVersion) > 0 {
		client.SetAPIVersion(c.APIVersion)
	} else {
		client.SetBaseURL(DefaultLinodeURL)
	}

	if c.EventPollMilliseconds != 0 {
		client.SetPollDelay(time.Duration(c.EventPollMilliseconds))
	}
	if c.MinRetryDelayMilliseconds != 0 {
		client.SetRetryWaitTime(time.Duration(c.MinRetryDelayMilliseconds) * time.Millisecond)
	}
	if c.MaxRetryDelayMilliseconds != 0 {
		client.SetRetryMaxWaitTime(time.Duration(c.MaxRetryDelayMilliseconds) * time.Millisecond)
	}

	return client
}

func terraformUserAgent(version string) string {
	ua := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s",
		version, meta.SDKVersionString())

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
