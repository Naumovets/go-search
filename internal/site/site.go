package site

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/logger/sl"
	"golang.org/x/net/html"
)

type TypeURL int

const (
	FullURL = iota
	RelativeURL
)

type Siter interface {
	NewSite() (*Site, error)
	NewChildSite() (*Site, error)
	newSiteRelative() (*Site, error)
	baseURL() string
	typeURL() TypeURL
	CompleteURL() string
	Analys() (string, []Site, error)
}

type Site struct {
	Id       int
	URL      string
	BasedURL string
	Type     TypeURL
	Content  string
}

func NewSite(URL string) (*Site, error) {

	isNotValid, err := isNotValidURL(URL)

	if err != nil {
		return nil, err
	}

	if isNotValid {
		return nil, fmt.Errorf("url is not valid: %s", URL)
	}

	w := Site{
		URL: URL,
	}

	if w.typeURL() == RelativeURL {
		return nil, fmt.Errorf("err: URL is not absolute: %s", URL)
	}
	w.baseURL()
	return &w, nil
}

func NewChildSite(URL string, parent Site) (*Site, error) {

	isNotValid, err := isNotValidURL(URL)

	if err != nil {
		return nil, err
	}

	if isNotValid {
		return nil, fmt.Errorf("url is not valid: %s", URL)
	}

	w := Site{
		URL: URL,
	}

	if w.typeURL() == RelativeURL {
		return newSiteRelative(URL, parent)
	}

	return NewSite(URL)
}

func newSiteRelative(URL string, parent Site) (*Site, error) {
	w := Site{
		URL: URL,
	}

	if w.typeURL() == FullURL {
		return nil, fmt.Errorf("err: URL isn't relative: %s", URL)
	}

	w.BasedURL = parent.URL

	return &w, nil

}

func (w *Site) baseURL() string {
	w.BasedURL = GetBaseURL(w.URL)
	return w.BasedURL
}

func (w *Site) typeURL() int {
	re, _ := regexp.Compile(`^((http|https)://|//)`)
	res := re.MatchString(w.URL)
	if res {
		w.Type = FullURL
		return FullURL
	}

	w.Type = RelativeURL
	return RelativeURL
}

func (w *Site) CompleteURL() (string, error) {
	if w.Type == FullURL {
		re, _ := regexp.Compile(`^//.+`)
		res := re.MatchString(w.URL)
		if res {
			return strings.Split(w.URL, "//")[1], nil
		}
		return w.URL, nil
	}

	ok, err := regexp.MatchString(`^/\w+`, w.URL)

	if err != nil {
		return "", fmt.Errorf("err: %w", err)
	}

	newURL := make([]string, 0)
	var domen string

	if ok {
		domen = GetBaseURL(w.BasedURL)
	} else {
		sliceBasedURL := make([]string, 0)

		for _, item := range strings.Split(w.BasedURL, "/") {
			if item != "" {
				sliceBasedURL = append(sliceBasedURL, item)
			}
		}

		var relativeURL []string
		if sliceBasedURL[0] == "http:" || sliceBasedURL[0] == "https:" {
			relativeURL = sliceBasedURL[2:]
			domen = strings.Join(sliceBasedURL[:2], "//")
		} else {
			domen = sliceBasedURL[0]
			relativeURL = sliceBasedURL[1:]
		}
		newURL = append(newURL, relativeURL...)
	}

	clearUrlSlice := make([]string, 0)

	for _, item := range strings.Split(w.URL, "/") {
		if item != "" {
			clearUrlSlice = append(clearUrlSlice, item)
		}
	}
	for _, item := range clearUrlSlice {
		if item == ".." {
			if len(newURL) > 0 {
				newURL = newURL[:len(newURL)-1]
			} else {
				return "", fmt.Errorf("err: relative url is failed: %v", newURL)
			}
		} else {
			newURL = append(newURL, item)
		}
	}

	resultURL := domen + "/" + strings.Join(newURL, "/")

	return resultURL, nil

}

func (site *Site) Analys() ([]Site, error) {
	resp, err := http.Get(site.URL)

	links := make([]string, 0)
	sites := make([]Site, 0)

	if err != nil {
		return nil, fmt.Errorf("err: %s", err)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("err: %s", err)
	}

	results := parse(doc, &links, false)

	for _, link := range links {
		newSite, err := NewChildSite(link, *site)

		if err != nil {
			log.Debug("Failed to add new site", sl.Err(err))
		} else {
			sites = append(sites, *newSite)
		}

	}

	site.Content = removeExtraSpaces(strings.Join(results, " "))

	return sites, nil
}
