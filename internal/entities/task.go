package entities

type Task struct {
	Id      int    `db:"id"`
	SiteURL string `db:"url"`
	Status  int    `db:"task_status"`
}
