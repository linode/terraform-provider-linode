//go:build tools
// +build tools

package main

//go:generate go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install golang.org/x/tools/cmd/goimports

import (
	// side effect imports used to version go tools
	// see: https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#tools-as-dependencies
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlint"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"
)
