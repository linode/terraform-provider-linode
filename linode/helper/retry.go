package helper

import (
	"log"
	"net/url"
	"regexp"

	"github.com/go-resty/resty/v2"
)

// Workaround for intermittent 5xx errors when retrieving a database from the API
func Database502Retry() func(response *resty.Response, err error) bool {
	databaseGetRegex, err := regexp.Compile("[A-Za-z0-9]+/databases/[a-z]+/instances/[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return func(response *resty.Response, _ error) bool {
		if response.StatusCode() != 502 || response.Request == nil {
			return false
		}

		requestURL, err := url.ParseRequestURI(response.Request.URL)
		if err != nil {
			log.Printf("[WARN] failed to parse request URL: %s", err)
			return false
		}

		// Check whether the string matches
		return databaseGetRegex.MatchString(requestURL.Path)
	}
}

func LinodeInstance500Retry() func(response *resty.Response, err error) bool {
	linodeGetRegex, err := regexp.Compile("linode/instances/[0-9]+/ips+")
	if err != nil {
		log.Fatal(err)
	}
	return func(response *resty.Response, _ error) bool {
		if response.StatusCode() != 500 || response.Request == nil {
			return false
		}

		requestURL, err := url.ParseRequestURI(response.Request.URL)
		if err != nil {
			log.Printf("[WARN] failed to parse request URL: %s", err)
			return false
		}

		// Check whether the string matches
		return linodeGetRegex.MatchString(requestURL.Path)
	}
}
