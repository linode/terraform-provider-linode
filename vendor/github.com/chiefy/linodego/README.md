# linodego

[![Build Status](https://travis-ci.org/chiefy/linodego.svg?branch=master)](https://travis-ci.org/chiefy/linodego)
[![GoDoc](https://godoc.org/github.com/chiefy/linodego?status.svg)](https://godoc.org/github.com/chiefy/linodego)

Go client for [Linode REST v4 API](https://developers.linode.com/v4/introduction)

## Installation

```
$ go get -u github.com/chiefy/linodego
```

## API Support

** Note: Currently in work-in-progress.  Things will change and break until we release a tagged version. **

Check [API_SUPPORT.md](API_SUPPORT.md) for current support of the Linode `v4` API endpoints.


## Documentation

See [godoc](https://godoc.org/github.com/chiefy/linodego) for a complete reference.

The API generally follows the naming patterns prescribed in the [OpenAPIv3 document for Linode APIv4](https://developers.linode.com/api/v4).

Deviations in naming have been made to avoid using "Linode" and "Instance" redundantly or inconsistently.

A brief summary of the features offered in this API client are shown here.

## Examples

### General Usage

```go
package main

import (
  "fmt"
  "log"
  "os"

  "github.com/chiefy/linodego"
)

func main() {
  apiKey, ok := os.LookupEnv("LINODE_TOKEN")
  if !ok {
    log.Fatal("Could not find LINODE_TOKEN, please assert it is set.")
  }
  linodeClient, err := linodego.NewClient(apiKey)
  if err != nil {
    log.Fatal(err)
  }
  linodeClient.SetDebug(true)
  res, err := linodeClient.GetInstance(4090913)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%v", res)

}
```

### Pagination
#### Auto-Pagination Requests

```go
kernels, err := linodego.ListKernels(nil)
// len(kernels) == 218
```

#### Single Page

```go
opts := NewListOptions(2,"")
// or opts := ListOptions{PageOptions: &PageOptions: {Page: 2 }}
kernels, err := linodego.ListKernels(opts)
// len(kernels) == 100
// opts.Results == 218
```

### Filtering

```go
opts := ListOptions{Filter: "{\"mine\":true}"}
// or opts := NewListOptions(0, "{\"mine\":true}")
stackscripts, err := linodego.ListStackscripts(opts)
```

### Error Handling
#### Getting Single Entities

```go
linode, err := linodego.GetLinode(555) // any Linode ID that does not exist or is not yours
// linode == nil: true
// err.Error() == "[404] Not Found"
// err.Code == "404"
// err.Message == "Not Found"
```

#### Lists

For lists, the list is still returned as `[]`, but `err` works the same way as on the `Get` request.

```go
linodes, err := linodego.ListLinodes(NewListOptions(0, "{\"foo\":bar}"))
// linodes == []
// err.Error() == "[400] [X-Filter] Cannot filter on foo"
```

Otherwise sane requests beyond the last page do not trigger an error, just an empty result:

```go
linodes, err := linodego.ListLinodes(NewListOptions(9999, ""))
// linodes == []
// err = nil
```

### Writes

When performing a `POST` or `PUT` request, multiple field related errors will be returned as a single error, currently like:

```go
// err.Error() == "[400] [field1] the problem reported by the API
// [field2] the problem reported by the API
// [field3] the problem reported by the API"
```

## Discussion / Help

Join us at [#linodego](https://gophers.slack.com/messages/CAG93EB2S) on the [gophers slack](https://gophers.slack.com)

## License

[MIT License](LICENSE)
