package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/version"
)

const UAEnvVar = "TF_APPEND_USER_AGENT"

// DefaultLinodeURL is the Linode APIv4 URL to use.
const DefaultLinodeURL = "https://api.linode.com"

type ProviderMeta struct {
	ProviderUserAgent string
	Client            linodego.Client
	Config            *Config
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
	SkipImplicitReboots          bool
	DisableInternalCache         bool
	MinRetryDelayMilliseconds    int
	MaxRetryDelayMilliseconds    int
	EventPollMilliseconds        int
	LKEEventPollMilliseconds     int
	LKENodeReadyPollMilliseconds int

	ObjAccessKey   string
	ObjSecretKey   string
	ObjUseTempKeys bool
}

// Client returns a fully initialized Linode client.
func (c *Config) Client(ctx context.Context, meta *ProviderMeta) (*linodego.Client, error) {
	loggingTransport := NewAPILoggerTransport(
		logging.NewSubsystemLoggingHTTPTransport(
			APILoggerSubsystem,
			http.DefaultTransport,
		),
	)

	oauth2Client := &http.Client{
		Transport: loggingTransport,
	}

	client := linodego.NewClient(oauth2Client)

	client.SetBaseURL(DefaultLinodeURL)

	// Load the config file if it exists
	if _, err := os.Stat(c.ConfigPath); err == nil {
		tflog.Info(ctx, "Using Linode profile", map[string]any{
			"config_path": c.ConfigPath,
		})
		err = client.LoadConfig(&linodego.LoadConfigOptions{
			Path:    c.ConfigPath,
			Profile: c.ConfigProfile,
		})
		if err != nil {
			return nil, err
		}
	} else {
		tflog.Info(ctx, "Linode config does not exist, skipping..")
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
		client.SetPollDelay(time.Duration(c.EventPollMilliseconds) * time.Millisecond)
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
	if meta != nil {
		meta.ProviderUserAgent = userAgent
	}

	ApplyAllRetryConditions(&client)

	// We always want to disable resty debugging in favor
	// of Terraform transport debugging.
	client.SetDebug(false)

	return &client, nil
}

func terraformUserAgent(version string) string {
	ua := fmt.Sprintf(
		"HashiCorp Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s",
		version, GetSDKv2Version(),
	)

	if add := os.Getenv(UAEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
