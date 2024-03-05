package helper

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func getDepVersion(depPath string) string {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		fmt.Println("No build information available")
		return ""
	}

	for _, dep := range info.Deps {
		if dep.Path == depPath {
			return strings.TrimLeft(dep.Version, "v")
		}
	}

	return ""
}

func GetSDKv2Version() string {
	return getDepVersion("github.com/hashicorp/terraform-plugin-sdk/v2")
}

func GetFrameworkVersion() string {
	return getDepVersion("github.com/hashicorp/terraform-plugin-framework")
}
