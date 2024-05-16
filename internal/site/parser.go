package site

import (
	"golang.org/x/net/html"
)

func parse(el *html.Node, sites *[]string, isText bool) []string {

	text := make([]string, 0)
	if el.Type == html.TextNode && isText {
		text = append(text, el.Data)
	}

	for c := el.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.ElementNode && c.Data == "a" {
			for _, item := range c.Attr {
				if item.Key == "href" && item.Val != "" {
					*sites = append(*sites, item.Val)
				}
			}
		}

		text = append(text, parse(c, sites, isTextTag(el))...)
	}

	return text

}
