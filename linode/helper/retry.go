package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
)

// Workaround for intermittent 5xx errors when retrieving a database from the API
func Database502Retry() func(response *http.Response, err error) bool {
	databaseGetRegex, err := regexp.Compile("[A-Za-z0-9]+/databases/[a-z]+/instances/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, databaseGetRegex)
}

func LinodeInstance500Retry() func(response *http.Response, err error) bool {
	linodeGetRegex, err := regexp.Compile("linode/instances/[0-9]+/ips+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, linodeGetRegex)
}

// ImageUpload500Retry for [500] error when uploading an image
func ImageUpload500Retry() func(response *http.Response, err error) bool {
	ImageUpload, err := regexp.Compile("images/upload")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, ImageUpload)
}

// OBJKeyCreate500Retry for [500] error when creating an Object Storage Key
func OBJKeyCreate500Retry() func(response *http.Response, err error) bool {
	OBJKeyCreate, err := regexp.Compile("object-storage/keys")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJKeyCreate)
}

// OBJKeyDelete500Retry for [500] error when deleting an Object Storage Key
func OBJKeyDelete500Retry() func(response *http.Response, err error) bool {
	OBJKeyDelete, err := regexp.Compile("object-storage/keys/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJKeyDelete)
}

// OBJBucketCreate500Retry for [500] error when creating an Object Storage Bucket
func OBJBucketCreate500Retry() func(response *http.Response, err error) bool {
	OBJBucketCreate, err := regexp.Compile("object-storage/buckets")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJBucketCreate)
}

// OBJBucketDelete500Retry for [500] error when deleting an Object Storage Bucket
func OBJBucketDelete500Retry() func(response *http.Response, err error) bool {
	OBJBucketDelete, err := regexp.Compile("object-storage/buckets/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJBucketDelete)
}

// InstanceDiskCreateBusyRetry retries disk creation when the Linode instance is busy (still provisioning).
// This handles transient errors that occur when attempting to create a disk immediately after instance creation:
// - "Linode busy." error (400 response with specific message)
// - Network/transport errors with [002] error code (e.g., "unexpected EOF", "failed to decode response body")
func InstanceDiskCreateBusyRetry() func(response *http.Response, err error) bool {
	diskCreatePath, compileErr := regexp.Compile("linode/instances/[0-9]+/disks$")
	if compileErr != nil {
		log.Fatal(compileErr)
	}

	return func(response *http.Response, err error) bool {
		// Extract context from request if available for better log correlation
		ctx := context.Background()
		if response != nil && response.Request != nil {
			ctx = response.Request.Context()
		}

		// Check if this is a [002] network/transport error on disk creation
		// This includes "unexpected EOF" and "failed to decode response body"
		if err != nil && strings.Contains(err.Error(), "[002]") &&
			(strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "failed to decode response body")) {
			tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: Checking [002] network error", map[string]any{
				"error": err.Error(),
			})
			// We need to verify this is a disk creation request
			if response != nil && response.Request != nil && response.Request.URL != nil {
				if response.Request.Method != http.MethodPost {
					return false
				}
				if diskCreatePath.MatchString(response.Request.URL.Path) {
					tflog.Debug(ctx, "Retrying disk creation due to [002] network error", map[string]any{
						"error": err.Error(),
						"path":  response.Request.URL.Path,
					})
					return true
				}
				tflog.Debug(ctx, "[002] network error on non-disk-creation endpoint, not retrying", map[string]any{
					"error": err.Error(),
					"path":  response.Request.URL.Path,
				})
			}
			return false
		}

		// Check for other errors (non-[002])
		if err != nil {
			tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: Non-[002] error, not retrying", map[string]any{
				"error": err.Error(),
			})
			return false
		}

		// Check for "Linode busy." response
		if response == nil || response.Request == nil || response.Request.URL == nil {
			tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: Nil response or request, not retrying", map[string]any{})
			return false
		}

		if response.StatusCode != 400 {
			return false
		}

		// Only retry disk creation requests (POST) on the disk collection endpoint
		if response.Request.Method != http.MethodPost {
			return false
		}

		// Only retry if the path matches disk creation
		if !diskCreatePath.MatchString(response.Request.URL.Path) {
			tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: 400 error on non-disk-creation endpoint, not retrying", map[string]any{
				"status_code": response.StatusCode,
				"path":        response.Request.URL.Path,
			})
			return false
		}

		// Check if the error message contains "Linode busy."
		// Need to read and restore the body since it can only be read once
		if response.Body == nil {
			tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: 400 response without body, not retrying", map[string]any{
				"path": response.Request.URL.Path,
			})
			return false
		}
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			tflog.Warn(ctx, "Failed to read response body", map[string]any{
				"error": err.Error(),
			})
			return false
		}
		response.Body.Close()
		response.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		bodyStr := string(bodyBytes)
		isLinodeBusy := strings.Contains(bodyStr, "Linode busy.")

		tflog.Debug(ctx, "InstanceDiskCreateBusyRetry: Checked 400 response body", map[string]any{
			"path":           response.Request.URL.Path,
			"is_linode_busy": isLinodeBusy,
			"will_retry":     isLinodeBusy,
		})

		return isLinodeBusy
	}
}

func GenericRetryCondition(statusCode int, pathPattern *regexp.Regexp) func(response *http.Response, err error) bool {
	return func(response *http.Response, _ error) bool {
		if response == nil ||
			response.Request == nil ||
			response.Request.URL == nil ||
			response.StatusCode != statusCode {
			return false
		}

		return pathPattern.MatchString(response.Request.URL.Path)
	}
}

func ApplyAllRetryConditions(client *linodego.Client) {
	client.AddRetryCondition(Database502Retry())
	client.AddRetryCondition(LinodeInstance500Retry())
	client.AddRetryCondition(ImageUpload500Retry())
	client.AddRetryCondition(OBJKeyCreate500Retry())
	client.AddRetryCondition(OBJKeyDelete500Retry())
	client.AddRetryCondition(OBJBucketCreate500Retry())
	client.AddRetryCondition(OBJBucketDelete500Retry())
	client.AddRetryCondition(InstanceDiskCreateBusyRetry())
}

// WithRetries runs the given retryFunc at most maxRetries times
// every interval if the returned bool is true.
func WithRetries(
	ctx context.Context,
	maxRetries int,
	interval time.Duration,
	retryFunc func() (bool, error),
) error {
	var lastError error

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for currentAttempt := range maxRetries {
		canRetry, err := retryFunc()
		if err == nil {
			// Success!
			return nil
		}

		lastError = err

		if !canRetry {
			return fmt.Errorf("got non-retryable error (attempt %d): %w", currentAttempt, err)
		}

		tflog.Warn(
			ctx,
			"Retrying failed operation",
			map[string]any{
				"attempt":      currentAttempt,
				"max_attempts": maxRetries,
				"error":        err.Error(),
			},
		)

		<-ticker.C
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastError)
}
