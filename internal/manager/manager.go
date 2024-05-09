package manager

// Пакет manager управляет поисковыми роботами

import (
	"fmt"
	"log/slog"

	"github.com/Naumovets/go-search/internal/collections/queue"
	"github.com/Naumovets/go-search/internal/entities"
	log "github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/logger/sl"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	"github.com/Naumovets/go-search/internal/robot"
	"github.com/Naumovets/go-search/internal/site"
)

type Manager struct {
	tasks      *queue.Queue[entities.Task]
	repository *tasks.Repository
}

func NewManager(rep *tasks.Repository) *Manager {
	return &Manager{
		tasks:      queue.NewQueue[entities.Task](),
		repository: rep,
	}
}

func (m *Manager) Start(count_tasks int) {
	for i := 0; i < 100; i++ {
		m.work(count_tasks)
	}

}

func (m *Manager) work(count_tasks int) {
	tasks, err := m.repository.GetActualLimitTasks(count_tasks)
	if err != nil {
		log.Debug("Failed to get actual tasks", sl.Err(err))
		return
	}

	if len(tasks) == 0 {
		log.Info("No tasks found")
		return
	}

	sites := make([]site.Site, 0)
	for id := range tasks {
		site, err := site.NewSite(tasks[id].SiteURL)
		site.Id = tasks[id].Id

		if err != nil {
			log.Debug("Failed to create new site", sl.Err(err))
			continue
		}

		sites = append(sites, *site)
	}

	if len(sites) == 0 {
		log.Debug(fmt.Sprintf("Count sites: %d", len(sites)))
		return
	}
	robot := robot.NewRobot(m.repository)
	robot.AddList(sites)
	new_sites, success_ids := robot.Work()

	if len(new_sites) > 0 {
		err := m.repository.AddTask(new_sites)
		if err != nil {
			log.Debug("Failed to add task", sl.Err(err))
			return
		}
	}

	if len(success_ids) > 0 {
		err = m.repository.CompleteTasks(success_ids)
		if err != nil {
			log.Debug("Failed to complete tasks", sl.Err(err))
			return
		}

		log.Info("Success done tasks", slog.String("array", fmt.Sprintf("%v", success_ids)))
	}
}
