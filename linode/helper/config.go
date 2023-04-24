package helper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/version"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

// DefaultLinodeURL is the Linode APIv4 URL to use.
const DefaultLinodeURL = "https://api.linode.com"

type FrameworkProviderMeta struct {
	Client *linodego.Client
}

type ProviderMeta struct {
	Client linodego.Client
	Config *Config
}

// Config represents the Linode provider configuration.
type Config struct {
	AccessToken string
	APIURL      string
	APIVersion  string
	UAPrefix    string

	ConfigPath    string
	ConfigProfile string

	TerraformVersion string

	SkipInstanceReadyPoll        bool
	SkipInstanceDeletePoll       bool
	DisableInternalCache         bool
	MinRetryDelayMilliseconds    int
	MaxRetryDelayMilliseconds    int
	EventPollMilliseconds        int
	LKEEventPollMilliseconds     int
	LKENodeReadyPollMilliseconds int
}

// Client returns a fully initialized Linode client.
func (c *Config) Client() (*linodego.Client, error) {
	loggingTransport := logging.NewTransport("Linode", http.DefaultTransport)

	oauth2Client := &http.Client{
		Transport: loggingTransport,
	}

	client := linodego.NewClient(oauth2Client)

	client.SetBaseURL(DefaultLinodeURL)

	// Load the config file if it exists
	if _, err := os.Stat(c.ConfigPath); err == nil {
		log.Println("[INFO] Using Linode profile: ", c.ConfigPath)
		err = client.LoadConfig(&linodego.LoadConfigOptions{
			Path:    c.ConfigPath,
			Profile: c.ConfigProfile,
		})
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("[INFO] Linode config does not exist, skipping..")
	}

	// Overrides
	if c.AccessToken != "" {
		client.SetToken(c.AccessToken)
	}

	if c.APIURL != "" {
		client.SetBaseURL(c.APIURL)
	}

	if len(c.APIVersion) > 0 {
		client.SetAPIVersion(c.APIVersion)
	}

	client.UseCache(!c.DisableInternalCache)

	if c.EventPollMilliseconds != 0 {
		client.SetPollDelay(time.Duration(c.EventPollMilliseconds))
	}
	if c.MinRetryDelayMilliseconds != 0 {
		client.SetRetryWaitTime(time.Duration(c.MinRetryDelayMilliseconds) * time.Millisecond)
	}
	if c.MaxRetryDelayMilliseconds != 0 {
		client.SetRetryMaxWaitTime(time.Duration(c.MaxRetryDelayMilliseconds) * time.Millisecond)
	}

	tfUserAgent := terraformUserAgent(c.TerraformVersion)
	userAgent := strings.TrimSpace(fmt.Sprintf("%s terraform-provider-linode/%s",
		tfUserAgent, version.ProviderVersion))
	if c.UAPrefix != "" {
		userAgent = c.UAPrefix + " " + userAgent
	}
	client.SetUserAgent(userAgent)

	// Workaround for intermittent 5xx errors when retrieving a database from the API
	databaseGetRegex, err := regexp.Compile("[A-Za-z0-9]+/databases/[a-z]+/instances/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	client.AddRetryCondition(func(response *resty.Response, err error) bool {
		if response.StatusCode() != 502 || response.Request == nil {
			return false
		}

		requestURL, err := url.ParseRequestURI(response.Request.URL)
		if err != nil {
			log.Printf("[WARN] failed to parse request URL: %s", err)
			return false
		}

		// Check whether the string matches
		return databaseGetRegex.MatchString(requestURL.Path)
	})

	return &client, nil
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
