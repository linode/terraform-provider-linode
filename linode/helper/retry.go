package helper

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/linode/linodego"
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

func GenericRetryCondition(statusCode int, pathPattern *regexp.Regexp) func(response *http.Response, err error) bool {
	return func(response *http.Response, _ error) bool {
		if response == nil || response.StatusCode != statusCode || response.Request == nil {
			return false
		}

		requestURL, err := url.ParseRequestURI(response.Request.RequestURI)
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
}
