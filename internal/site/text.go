package site

import (
	"regexp"

	"golang.org/x/net/html"
)

func isTextTag(c *html.Node) bool {
	htmlTags := map[string]interface{}{
		"a":     interface{}(nil),
		"p":     interface{}(nil),
		"span":  interface{}(nil),
		"h1":    interface{}(nil),
		"h2":    interface{}(nil),
		"h3":    interface{}(nil),
		"h4":    interface{}(nil),
		"h5":    interface{}(nil),
		"h6":    interface{}(nil),
		"li":    interface{}(nil),
		"title": interface{}(nil),
	}

	if _, ok := htmlTags[c.Data]; ok && c.Type == html.ElementNode {
		return true
	}
	return false
}

func removeExtraSpaces(text string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(text, " ")
}

func IsRussianText(text string) bool {
	russianCount := 0
	englishCount := 0

	for _, char := range text {
		if (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я') {
			russianCount++
		} else if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			englishCount++
		}
	}

	return russianCount > englishCount
}
