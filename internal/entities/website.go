package entities

type Website struct {
	URL     string `db:"url" json:"url"`
	Content string `db:"content" json:"content"`
	Title   string `db:"title" json:"title"`
	// ContentTSV tsvector.TSVector `db:"content_tsv"`
	Language string  `db:"language" json:"language,omitempty"`
	Rank     float32 `db:"rank" json:"rank,omitempty"`
}
