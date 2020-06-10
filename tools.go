// +build tools

package main

import (
	// side effect imports used to version go tools
	// see: https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#tools-as-dependencies
	_ "github.com/Charliekenney23/tf-changelog-validator/cmd/tf-changelog-validator"
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
