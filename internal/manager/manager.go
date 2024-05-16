package manager

// Пакет manager управляет поисковыми роботами

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"

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
	mu         sync.Mutex
}

func NewManager(rep *tasks.Repository) *Manager {
	return &Manager{
		tasks:      queue.NewQueue[entities.Task](),
		repository: rep,
	}
}

func (m *Manager) Start(count_tasks int) {

	var wg sync.WaitGroup

	maxProcs := 5

	wg.Add(maxProcs)
	runtime.GOMAXPROCS(maxProcs)

	for i := 0; i < maxProcs; i++ {
		i := i
		go func() {
			defer wg.Done()
			log.Info(fmt.Sprintf("work #%d", i))
			m.work(count_tasks)
		}()
	}

	wg.Wait()
}

func (m *Manager) work(count_tasks int) {

	m.mu.Lock()
	tasks, err := m.repository.GetActualLimitTasks(count_tasks)
	m.mu.Unlock()
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

		if tasks[id].SiteURL == "" {
			log.Info(fmt.Sprintf("url: %s id: %d status: %d", tasks[id].SiteURL, tasks[id].Id, tasks[id].Status))
		}

		site, err := site.NewSite(tasks[id].SiteURL)

		if err != nil {
			log.Debug(fmt.Sprintf("Failed to create new site: %s", tasks[id].SiteURL), sl.Err(err))
			continue
		}
		site.Id = tasks[id].Id

		sites = append(sites, *site)
	}

	if len(sites) == 0 {
		log.Debug(fmt.Sprintf("Count sites: %d", len(sites)))
		return
	}
	robot := robot.NewRobot(m.repository)
	robot.AddList(sites)
	new_sites, success_ids := robot.Work()

	log.Info(fmt.Sprintf("c. sites: %d c. success_ids: %d", len(new_sites), len(success_ids)))

	if len(new_sites) > 0 {
		err := m.repository.AddTask(new_sites)
		if err != nil {
			log.Debug("Failed to add task", sl.Err(err))
		}
	}

	if len(success_ids) > 0 {
		err = m.repository.CompleteTasks(success_ids)
		if err != nil {
			log.Debug("Failed to complete tasks", sl.Err(err))
			return
		}

		log.Info("Success done tasks:", slog.String("count", fmt.Sprintf("%v", len(success_ids))))
		return
	}

	log.Info("No success done tasks")
}
