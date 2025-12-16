package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linode/linodego"
)

// HTTPClientModifier is the signature for functions used to modify an HTTP client before use.
type HTTPClientModifier func(client *http.Client) error

// AddRootCAToTransport applies the cert file at the given path to the given *http.Transport
func AddRootCAToTransport(cert string, transport *http.Transport) error {
	certData, err := os.ReadFile(filepath.Clean(cert))
	if err != nil {
		return fmt.Errorf("failed to read cert file %s: %w", cert, err)
	}

	tlsConfig := transport.TLSClientConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if tlsConfig.RootCAs == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
	}

	tlsConfig.RootCAs.AppendCertsFromPEM(certData)

	return nil
}

// NotFoundDefault wraps the given inner function, returning the given fallback
// value if it returns a linodego 404 error.
func NotFoundDefault[T any](inner func() (T, error), fallback T) (T, error) {
	result, err := inner()
	if linodego.IsNotFound(err) {
		return fallback, nil
	}

	return result, err
}

// NotFoundDefaultSlice functions the same as NotFoundDefault but uses an empty slice
// as a fallback.
func NotFoundDefaultSlice[T any](inner func() ([]T, error)) ([]T, error) {
	return NotFoundDefault(inner, nil)
}
