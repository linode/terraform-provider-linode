package helper

import (
	"fmt"
	"runtime/debug"
)

func getDepVersion(depPath string) string {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		fmt.Println("No build information available")
		return ""
	}

	for _, dep := range info.Deps {
		if dep.Path == depPath {
			return dep.Version
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
