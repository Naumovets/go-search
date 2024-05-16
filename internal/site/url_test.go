package site

import (
	"fmt"
	"testing"
)

func TestGetBaseUrl(t *testing.T) {
	urls := []string{
		"https://www.example.com/gdfdfdfd/fdfdfdfdsfd/fdsfsfds",
		"http://subdomain.example.com",
		"//another.example.com",
		"www.example.co.uk",
		"sub1.sub2.example.com.ru",
		"coursera.online",
		"пример.сайта.рф.ком.ру/тест1/тест2/тест3",
	}

	for _, url := range urls {
		res := GetBaseURL(url)

		fmt.Println(res)

	}
}
