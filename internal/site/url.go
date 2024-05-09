package site

import (
	"fmt"
	"regexp"
	"strings"
)

func GetBaseURL(url string) string {
	re, _ := regexp.Compile(`^((http|https)://|//|)\w+[.]\w+`)
	res := re.FindString(url)

	ok, _ := regexp.MatchString(`^//\w+`, res)

	if ok {
		res = strings.Split(res, "//")[1]
	}

	return res
}

func isNotValidURL(url string) (bool, error) {

	if url == "" {
		return false, fmt.Errorf("err: URL must not be empty")
	}

	ok, err := regexp.MatchString(`^(mailto|tel):.+`, url)

	if err != nil {
		return false, err
	}

	return ok, nil
}
