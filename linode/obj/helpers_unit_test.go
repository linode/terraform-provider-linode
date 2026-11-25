//go:build unit

package obj

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

func TestCleanUpTempKeys_SuccessFirstAttempt(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v4/object-storage/keys/123", r.URL.Path)
		attempts.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	cleanUpTempKeysWithInterval(context.Background(), &client, 123, 3, 5*time.Millisecond)
	assert.Equal(t, int32(1), attempts.Load())
}

func TestCleanUpTempKeys_RetryOn500ThenSuccess(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v4/object-storage/keys/123", r.URL.Path)
		current := attempts.Add(1)
		if current == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors":[{"reason":"Internal server error"}]}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	cleanUpTempKeysWithInterval(context.Background(), &client, 123, 3, 5*time.Millisecond)
	assert.Equal(t, int32(2), attempts.Load())
}

func TestCleanUpTempKeys_AlreadyDeleted404(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v4/object-storage/keys/123", r.URL.Path)
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"reason":"Not found"}]}`))
	}))
	defer server.Close()

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	cleanUpTempKeysWithInterval(context.Background(), &client, 123, 3, 5*time.Millisecond)
	assert.Equal(t, int32(1), attempts.Load())
}

func TestCleanUpTempKeys_MaxRetriesExhausted(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v4/object-storage/keys/123", r.URL.Path)
		attempts.Add(1)
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"errors":[{"reason":"Bad gateway"}]}`))
	}))
	defer server.Close()

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	cleanUpTempKeysWithInterval(context.Background(), &client, 123, 3, 5*time.Millisecond)
	assert.Equal(t, int32(3), attempts.Load())
}
