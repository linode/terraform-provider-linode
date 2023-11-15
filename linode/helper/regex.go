package helper

import (
	"fmt"
	"regexp"
)

func StringToRegex(pattern string) (regExp *regexp.Regexp) {
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	return regExp
}
