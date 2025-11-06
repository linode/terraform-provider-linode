package acceptance

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"
)

func NewResponseOverrideClient(
	t *testing.T,
	shouldOverride func(response *http.Response) bool,
	override func(t *testing.T, responseBody map[string]any),
) *linodego.Client {
	client, err := GetTestClient()
	require.NoError(t, err)

	// TODO: Expose client URL through public interface
	rawURL := cmp.Or(os.Getenv("LINODE_URL"), linodego.APIHost)
	if !strings.Contains(rawURL, "://") {
		// Assume HTTPS
		rawURL = fmt.Sprintf("%s://%s", linodego.APIProto, rawURL)
	}

	parsedURL, err := url.Parse(rawURL)
	require.NoError(t, err)

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = parsedURL.Scheme
			req.URL.Host = parsedURL.Host
			req.Host = parsedURL.Host
		},
		ModifyResponse: func(resp *http.Response) error {
			if !shouldOverride(resp) {
				return nil
			}

			var body map[string]any
			err := json.NewDecoder(resp.Body).Decode(&body)
			require.NoError(t, err)

			override(t, body)

			var buf bytes.Buffer
			require.NoError(t, json.NewEncoder(&buf).Encode(body))

			require.NoError(t, resp.Body.Close())
			resp.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

			resp.ContentLength = int64(buf.Len())
			resp.Header.Set("Content-Length", strconv.Itoa(buf.Len()))

			return nil
		},
	}

	proxyServer := httptest.NewServer(proxy)
	t.Cleanup(proxyServer.Close)

	proxiedClient, err := client.UseURL(proxyServer.URL)
	require.NoError(t, err)

	return proxiedClient
}
