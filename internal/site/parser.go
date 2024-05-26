package site

import (
	"golang.org/x/net/html"
)

func getTitle(el html.Node, isText bool) []string {

	title := make([]string, 0)
	if el.Type == html.TextNode && isText {
		title = append(title, el.Data)
	}

	for c := el.FirstChild; c != nil; c = c.NextSibling {
		title = append(title, getTitle(*c, isTextTag(&el))...)
	}

	return title
}

func parse(el *html.Node, sites *[]string, isText bool) ([]string, []string) {

	text := make([]string, 0)
	title := make([]string, 0)
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
		} else if c.Type == html.ElementNode && c.Data == "title" {
			title = append(title, getTitle(*c, true)...)
		}

		newText, newTitle := parse(c, sites, isTextTag(el))
		text = append(text, newText...)
		title = append(title, newTitle...)
	}

	return text, title

}
