package helper

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// Workaround for intermittent 5xx errors when retrieving a database from the API
func Database502Retry() func(response *resty.Response, err error) bool {
	databaseGetRegex, err := regexp.Compile("[A-Za-z0-9]+/databases/[a-z]+/instances/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, databaseGetRegex)
}

func LinodeInstance500Retry() func(response *resty.Response, err error) bool {
	linodeGetRegex, err := regexp.Compile("linode/instances/[0-9]+/ips+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, linodeGetRegex)
}

// ImageUpload500Retry for [500] error when uploading an image
func ImageUpload500Retry() func(response *resty.Response, err error) bool {
	ImageUpload, err := regexp.Compile("images/upload")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, ImageUpload)
}

// OBJKeyCreate500Retry for [500] error when creating an Object Storage Key
func OBJKeyCreate500Retry() func(response *resty.Response, err error) bool {
	OBJKeyCreate, err := regexp.Compile("object-storage/keys")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJKeyCreate)
}

// OBJKeyDelete500Retry for [500] error when deleting an Object Storage Key
func OBJKeyDelete500Retry() func(response *resty.Response, err error) bool {
	OBJKeyDelete, err := regexp.Compile("object-storage/keys/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJKeyDelete)
}

// OBJBucketCreate500Retry for [500] error when creating an Object Storage Bucket
func OBJBucketCreate500Retry() func(response *resty.Response, err error) bool {
	OBJBucketCreate, err := regexp.Compile("object-storage/buckets")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJBucketCreate)
}

// OBJBucketDelete500Retry for [500] error when deleting an Object Storage Bucket
func OBJBucketDelete500Retry() func(response *resty.Response, err error) bool {
	OBJBucketDelete, err := regexp.Compile("object-storage/buckets/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return GenericRetryCondition(500, OBJBucketDelete)
}

// InstanceDiskCreateBusyRetry retries disk creation when the Linode instance is busy (still provisioning).
// This handles transient errors that occur when attempting to create a disk immediately after instance creation:
// - "Linode busy." error (400 response with specific message)
// - EOF errors (network/connection issues during provisioning)
func InstanceDiskCreateBusyRetry() func(response *resty.Response, err error) bool {
	diskCreatePath, compileErr := regexp.Compile("linode/instances/[0-9]+/disks$")
	if compileErr != nil {
		log.Fatal(compileErr)
	}

	return func(response *resty.Response, err error) bool {
		// Check if this is an EOF error on disk creation
		if err != nil && strings.Contains(err.Error(), "EOF") {
			// We need to verify this is a disk creation request
			if response != nil && response.Request != nil {
				requestURL, urlErr := url.ParseRequestURI(response.Request.URL)
				if urlErr == nil && diskCreatePath.MatchString(requestURL.Path) {
					log.Printf("[DEBUG] Retrying disk creation due to EOF error: %s", err.Error())
					return true
				}
			}
			return false
		}

		// Check for "Linode busy." response
		if response == nil || response.Request == nil {
			return false
		}

		if response.StatusCode() != 400 {
			return false
		}

		requestURL, urlErr := url.ParseRequestURI(response.Request.URL)
		if urlErr != nil {
			log.Printf("[WARN] failed to parse request URL: %s", urlErr)
			return false
		}

		// Only retry if the path matches disk creation
		if !diskCreatePath.MatchString(requestURL.Path) {
			return false
		}

		// Check if the error message contains "Linode busy."
		bodyStr := string(response.Body())
		return strings.Contains(bodyStr, "Linode busy.")
	}
}

func GenericRetryCondition(statusCode int, pathPattern *regexp.Regexp) func(response *resty.Response, err error) bool {
	return func(response *resty.Response, _ error) bool {
		if response.StatusCode() != statusCode || response.Request == nil {
			return false
		}

		requestURL, err := url.ParseRequestURI(response.Request.URL)
		if err != nil {
			log.Printf("[WARN] failed to parse request URL: %s", err)
			return false
		}

		// Check whether the string matches
		return pathPattern.MatchString(requestURL.Path)
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
