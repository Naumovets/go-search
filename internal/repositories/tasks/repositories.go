package tasks

import (
	"fmt"
	"strings"

	"github.com/Naumovets/go-search/internal/entities"
	"github.com/Naumovets/go-search/internal/site"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddTask(sites []site.Site) error {
	valueQuery := make([]string, 0)
	valueArgs := make([]interface{}, 0) // Изменение: interface{} для универсальности

	i := 1
	for _, site := range sites {
		url, err := site.CompleteURL()
		if err != nil {
			// Обработка ошибки получения URL (например, логирование)
			continue // Пропустить этот сайт, если возникла ошибка
		}

		valueQuery = append(valueQuery, fmt.Sprintf("($%d)", i))
		valueArgs = append(valueArgs, url) // Добавляем url как interface{}
		i++
	}

	query := fmt.Sprintf(`
	  INSERT INTO task (url)
	  VALUES %s
	  ON CONFLICT (id) DO NOTHING;
	`, strings.Join(valueQuery, ", "))

	_, err := r.db.Exec(query, valueArgs...) // Использование оператора "..." для развертывания среза

	return err // Вернуть ошибку, если она есть
}

func (r *repository) GetActualLimitTasks(lim int) ([]entities.Task, error) {
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

func (r *repository) GetLimitTasks(lim int) ([]entities.Task, error) {
	if lim <= 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT id, url, task_status FROM task LIMIT %d`, lim)

	rawTask := make([]entities.Task, 0)
	err := r.db.Select(&rawTask, query)

	if err != nil {
		return nil, err
	}

	return rawTask, nil
}

func (r *repository) CompleteTasks(ids []int) error {
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
