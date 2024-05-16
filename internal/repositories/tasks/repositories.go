package tasks

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

func (r *Repository) AddTask(sites []site.Site) error {

	// TODO: to add tasks in part for avoiding pg exeption

	var err error

	for i := 0; i < (len(sites)/1000)+1; i++ {
		k := 1
		valueQuery := make([]string, 0)
		valueArgs := make([]interface{}, 0)
		for j := 1000 * i; j < min(len(sites), (i+1)*1000); j++ {
			url, err := sites[j].CompleteURL()
			if err != nil {
				log.Debug("Failed to complete url", sl.Err(err))
				continue
			}
			valueQuery = append(valueQuery, fmt.Sprintf("($%d)", k))
			valueArgs = append(valueArgs, url)
			k++
		}

		query := fmt.Sprintf(`
			INSERT INTO task (url)
			VALUES %s
			ON CONFLICT (url)
			DO NOTHING`,
			strings.Join(valueQuery, ", "))

		_, err = r.db.Exec(query, valueArgs...)
	}

	return err
}

func (r *Repository) GetActualLimitTasks(lim int) ([]entities.Task, error) {
	if lim <= 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
	WITH next_task AS (
		SELECT id FROM task
		WHERE task_status = 0
		LIMIT %d
		FOR UPDATE SKIP LOCKED
	)
	UPDATE task
	set
		task_status = 1
	FROM next_task
	WHERE task.id = next_task.id
	RETURNING task.id, task.url, task.task_status;
	`, lim)

	rawTask := make([]entities.Task, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return nil, err
	}

	return rawTask, nil
}

func (r *Repository) GetLimitTasks(lim int) ([]entities.Task, error) {
	if lim <= 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT id, url, task_status FROM task ORDER BY id LIMIT %d`, lim)

	rawTask := make([]entities.Task, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return nil, err
	}

	return rawTask, nil
}

func (r *Repository) ExistActualTask() (bool, error) {
	query := "SELECT id FROM task WHERE task_status = 0 LIMIT 1"
	rawTask := make([]entities.Task, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return false, err
	}

	if len(rawTask) == 0 {
		return false, nil
	}

	return true, nil
}

func (r *Repository) CompleteTasks(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	valueQuery := make([]string, 0)
	valueArgs := make([]interface{}, 0)

	for i, id := range ids {
		valueQuery = append(valueQuery, fmt.Sprintf("$%d", i+1))
		valueArgs = append(valueArgs, id)
	}

	query := fmt.Sprintf(`UPDATE task
							SET
								task_status = 2
							WHERE id IN (%s);`, strings.Join(valueQuery, ", "))
	_, err := r.db.Exec(query, valueArgs...)

	if err != nil {
		return err
	}

	return nil

}
