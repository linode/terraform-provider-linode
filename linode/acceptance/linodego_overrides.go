package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type (
	shouldOverrideResponseFunc func(response *http.Response) bool
	responseOverrideFunc       func(responseBody map[string]any) error
)

type responseOverrideTransport struct {
	next             http.RoundTripper
	shouldOverride   shouldOverrideResponseFunc
	overrideResponse responseOverrideFunc
}

func (t *responseOverrideTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.next.RoundTrip(request)
	if err != nil {
		return response, err
	}

	if !t.shouldOverride(response) {
		return response, nil
	}

	var bodyData map[string]any
	if err := json.NewDecoder(response.Body).Decode(&bodyData); err != nil {
		return response, fmt.Errorf("failed to decode response body json: %w", err)
	}

	if err := t.overrideResponse(bodyData); err != nil {
		return response, fmt.Errorf("failed to override response body: %w", err)
	}

	var resultBuffer bytes.Buffer

	if err := json.NewEncoder(&resultBuffer).Encode(bodyData); err != nil {
		return response, fmt.Errorf("failed to encode modified response body: %w", err)
	}

	response.Body = io.NopCloser(bytes.NewReader(resultBuffer.Bytes()))

	return response, nil
}

func NewResponseOverrideClient(
	t *testing.T,
	shouldOverride shouldOverrideResponseFunc,
	overrideResponse responseOverrideFunc,
) *linodego.Client {
	t.Helper()

	return GetFrameworkTestClient(
		t,
		[]helper.HTTPClientModifier{
			func(client *http.Client) error {
				client.Transport = &responseOverrideTransport{
					next:             client.Transport,
					shouldOverride:   shouldOverride,
					overrideResponse: overrideResponse,
				}
				return nil
			},
		},
	)
}
