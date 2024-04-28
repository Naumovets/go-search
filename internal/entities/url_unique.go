package entities

type Task struct {
	Id      int    `db:"id"`
	SiteURL string `db:"url"`
	Status  string `db:"task_status"`
}
