package helper

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	APILoggerSubsystem = "linodego-requests"
	EnvAPILogger       = "TF_LOG_PROVIDER_LINODE_REQUESTS"
)

var APILogLevel = hclog.LevelFromString(os.Getenv(EnvAPILogger))

// APILoggerTransport injects a configured API request logger subsystem
// into the context of the current API request.
type APILoggerTransport struct {
	transport http.RoundTripper
}

// NewAPILoggerTransport is a RoundTripper used to inject
// a subsystem logger at request-time.
func NewAPILoggerTransport(transport http.RoundTripper) *APILoggerTransport {
	return &APILoggerTransport{
		transport: transport,
	}
}

// RoundTrip injects the API logger subsystem into the context
// of an API request. This allows us to configure the logger without
// creating a new logger in each implementation.
func (t *APILoggerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.transport.RoundTrip(r.WithContext(t.createAPILoggerSubsystem(r.Context())))
}

// createAPILoggerSubsystem creates an API logger subsystem
// for logging raw HTTP requests.
func (t *APILoggerTransport) createAPILoggerSubsystem(ctx context.Context) context.Context {
	targetLevel := APILogLevel

	// Disable the logger is no logger is defined
	if targetLevel == hclog.NoLevel {
		targetLevel = hclog.Off
	}

	ctx = tflog.NewSubsystem(ctx, APILoggerSubsystem, tflog.WithLevel(targetLevel))
	ctx = tflog.SubsystemMaskFieldValuesWithFieldKeys(ctx, APILoggerSubsystem, "Authorization")
	return ctx
}
