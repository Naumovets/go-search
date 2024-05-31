package entities

import "github.com/aymericbeaumet/go-tsvector"

type Website struct {
	URL        string            `db:"url"`
	Content    string            `db:"content"`
	Title      string            `db:"title"`
	ContentTSV tsvector.TSVector `db:"content_tsv"`
	Language   string            `db:"language"`
	Rank       float32           `db:"rank"`
}
