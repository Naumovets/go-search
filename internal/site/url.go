package site

import (
	"fmt"
	"regexp"
)

func GetBaseURL(url string) string {
	re, _ := regexp.Compile(`^((http|https)://|)\w+[.]\w+`)
	res := re.FindString(url)
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
