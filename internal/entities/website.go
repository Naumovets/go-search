package entities

type website struct {
	URL     string `db:"url"`
	Content string `db:"content"`
	Title   string `db:"title"`
}
