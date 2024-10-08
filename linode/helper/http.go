package helper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// AddRootCAToTransport applies the cert file at the given path to the given *http.Transport
func AddRootCAToTransport(cert string, transport *http.Transport) error {
	certData, err := os.ReadFile(cert)
	if err != nil {
		return fmt.Errorf("failed to read cert file %s: %w", cert, err)
	}

	tlsConfig := transport.TLSClientConfig
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	if tlsConfig.RootCAs == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
	}

	tlsConfig.RootCAs.AppendCertsFromPEM(certData)

	return nil
}
