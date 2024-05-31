package manager

import (
	"github.com/Naumovets/go-search/internal/site"
)

func getCompletedSites(sites []*site.Site, ids []int) []site.Site {
	result := make([]site.Site, len(ids))

	k := 0
	for _, site := range sites {
		if contains(ids, site.Id) {
			result[k] = *site
			k++
		}
	}
	return result
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
