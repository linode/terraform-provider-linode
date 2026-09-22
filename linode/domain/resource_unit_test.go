//go:build unit

package domain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) linodego.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)

	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	return client
}

func TestGetDomainWithRetries_RetriesOn404(t *testing.T) {
	var requests atomic.Int32

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if requests.Add(1) < 3 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errors": [{"reason": "Not found"}]}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": 1234, "domain": "example.org", "type": "master"}`))
	})

	domain, err := getDomainWithRetries(
		context.Background(), client, 1234,
		5*time.Second, time.Millisecond, 10*time.Millisecond,
	)

	require.NoError(t, err)
	require.NotNil(t, domain)
	assert.Equal(t, 1234, domain.ID)
	assert.Equal(t, "example.org", domain.Domain)
	assert.EqualValues(t, 3, requests.Load())
}

func TestGetDomainWithRetries_DeadlineExceeded(t *testing.T) {
	var requests atomic.Int32

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors": [{"reason": "Not found"}]}`))
	})

	domain, err := getDomainWithRetries(
		context.Background(), client, 1234,
		50*time.Millisecond, time.Millisecond, 10*time.Millisecond,
	)

	require.Error(t, err)
	assert.Nil(t, domain)
	assert.True(t, linodego.IsNotFound(err), "expected the last 404 error to be returned")
	assert.Greater(t, requests.Load(), int32(1), "expected at least one retry")
}

func TestGetDomainWithRetries_NoRetryOnNon404(t *testing.T) {
	var requests atomic.Int32

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors": [{"reason": "Bad request"}]}`))
	})

	domain, err := getDomainWithRetries(
		context.Background(), client, 1234,
		5*time.Second, time.Millisecond, 10*time.Millisecond,
	)

	require.Error(t, err)
	assert.Nil(t, domain)
	assert.False(t, linodego.IsNotFound(err))
	assert.EqualValues(t, 1, requests.Load())
}

func TestGetDomainWithRetries_ContextCancelled(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors": [{"reason": "Not found"}]}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	domain, err := getDomainWithRetries(
		ctx, client, 1234,
		5*time.Second, time.Millisecond, 10*time.Millisecond,
	)

	require.Error(t, err)
	assert.Nil(t, domain)
}
