package helper

import (
	"log"
	"net/url"
	"regexp"

	"github.com/go-resty/resty/v2"
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
}
