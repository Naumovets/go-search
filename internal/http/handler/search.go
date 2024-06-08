package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/repositories/db"
)

type Search struct {
	rep *db.Repository
}

func New(rep *db.Repository) *Search {
	return &Search{
		rep: rep,
	}
}

func (s *Search) Search(w http.ResponseWriter, r *http.Request) {

	strPage := strings.TrimSpace(r.URL.Query().Get("page"))
	var page int
	var err error

	if strPage == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(strPage)
	}

	if err != nil || page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("page must be positive interger, was got: '%s'", strPage)))
		return
	}

	strLim := strings.TrimSpace(r.URL.Query().Get("lim"))
	var lim int
	if strLim == "" {
		lim = 20
	} else {
		lim, err = strconv.Atoi(strLim)
	}

	if err != nil || lim < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("lim must be positive interger, was got: '%s'", strLim)))
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("query"))

	reg := regexp.MustCompile(`\s+`)
	query = reg.ReplaceAllString(query, " ")

	if query == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	websites, err := s.rep.Search(query, lim, page)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug(fmt.Sprintf("failed to search this query: %s", err))
		return
	}

	response, err := json.Marshal(websites)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug(fmt.Sprintf("failed to encode websites to json: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
