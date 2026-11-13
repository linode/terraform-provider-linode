//go:build unit

package helper

import (
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRootCAToTransport_TrustsCert(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	caPath := filepath.Join(t.TempDir(), "ca.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	require.NoError(t, os.WriteFile(caPath, certPEM, 0o600))

	// This mirrors how the provider builds its transport.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	require.NoError(t, AddRootCAToTransport(caPath, transport))

	resp, err := (&http.Client{Transport: transport}).Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAddRootCAToTransport_InvalidPEM(t *testing.T) {
	caPath := filepath.Join(t.TempDir(), "ca.pem")
	require.NoError(t, os.WriteFile(caPath, []byte("not a certificate"), 0o600))

	transport := http.DefaultTransport.(*http.Transport).Clone()
	require.ErrorContains(t, AddRootCAToTransport(caPath, transport), "failed to parse any PEM certificates")
}
