# go-linode

[![Build Status](https://travis-ci.org/chiefy/go-linode.svg?branch=master)](https://travis-ci.org/chiefy/go-linode)

Go client for [Linode REST v4 API](https://developers.linode.com/v4/introduction)

## Installation

```
$ go get -u github.com/chiefy/go-linode
```

## API Support

** Note: currently pagination is not supported. The response list will return the first page of responses only **

Check [API_SUPPORT.md](API_SUPPORT.md) for current support of the Linode `v4` API endpoints.


## Documentation

Current in progress.

## Example

```go
package main

import (
  "fmt"
  "log"
  "os"

  golinode "github.com/chiefy/go-linode"
)

func main() {
  apiKey, ok := os.LookupEnv("LINODE_API_KEY")
  if !ok {
    log.Fatal("Could not find LINODE_API_KEY, please assert it is set.")
  }
  linodeClient, err := golinode.NewClient(apiKey)
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

## Discussion / Help

Join us at `#go-linode` on the [gophers slack](https://gophers.slack.com)

## License

[MIT License](LICENSE)
