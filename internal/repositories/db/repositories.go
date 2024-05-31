package db

import (
	"fmt"
	"strings"

	"github.com/Naumovets/go-search/internal/entities"
	log "github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/logger/sl"
	"github.com/Naumovets/go-search/internal/site"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// TODO: func for adding website to db

func (r *Repository) AddSites(sites []site.Site) error {
	var err error

	for i := 0; i < (len(sites)/1000)+1; i++ {
		dbSites := make([]entities.Website, 0)
		for j := 1000 * i; j < min(len(sites), (i+1)*1000); j++ {
			url, err := sites[j].CompleteURL()
			content := sites[j].Content
			title := sites[j].Title
			if err != nil {
				log.Debug("Failed to complete url", sl.Err(err))
				continue
			}
			if site.IsRussianText(content) {
				dbSites = append(dbSites, entities.Website{
					URL:      url,
					Content:  content,
					Title:    title,
					Language: "russian",
				})
			} else {
				dbSites = append(dbSites, entities.Website{
					URL:      url,
					Content:  strings.TrimSpace(content),
					Title:    strings.TrimSpace(title),
					Language: "english",
				})
			}
		}

		query := `
		INSERT INTO website (url, content, title, content_tsv)
		VALUES (:url,:content,:title,to_tsvector(:language,:content))
		ON CONFLICT (url)
		DO NOTHING
		`
		_, err = r.db.NamedExec(query, dbSites)
	}

	return err
}

func (r *Repository) GetLimitSites(lim int) ([]entities.Website, error) {
	if lim <= 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT url, content, title, content_tsv FROM website LIMIT %d`, lim)

	rawTask := make([]entities.Website, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return nil, err
	}

	return rawTask, nil
}

func (r *Repository) Search(userQuery string, lim int, page int) ([]entities.Website, error) {
	if lim <= 0 {
		return nil, fmt.Errorf("lim must be 1 or more")
	}

	if page <= 0 {
		return nil, fmt.Errorf("page must be 1 or more")
	}
	var lang string
	if site.IsRussianText(userQuery) {
		lang = "russian"
	} else {
		lang = "english"
	}

	query := fmt.Sprintf(`
					SELECT url, content, title 
					FROM website 
					WHERE content_tsv @@ plainto_tsquery('%s', '%s')
					LIMIT %d 
					OFFSET %d`,
		lang,
		userQuery,
		lim,
		lim*(page-1))

	rawTask := make([]entities.Website, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return nil, err
	}

	return rawTask, nil
}
