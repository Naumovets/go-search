package manager

// Пакет manager управляет поисковыми роботами

import (
	"fmt"
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
	queue_tasks *queue.Queue[entities.Task]
	repository  *tasks.Repository
	mu          sync.Mutex
}

func NewManager(rep *tasks.Repository) *Manager {
	return &Manager{
		queue_tasks: queue.NewQueue[entities.Task](),
		repository:  rep,
	}
}

func (m *Manager) Start(count_tasks int) {

	for {
		// get tasks to queue
		tasks, err := m.repository.GetNoCompleteLimitTasks(count_tasks)
		log.Info(fmt.Sprintf("%d", len(tasks)))
		var actual_tasks []entities.Task

		if err != nil {
			log.Debug("Failed to get no complete tasks", sl.Err(err))
			actual_tasks, err = m.repository.GetActualLimitTasks(count_tasks)
			if err != nil {
				log.Debug("Failed to get tasks", sl.Err(err))
				return
			}
			tasks = actual_tasks

		} else if len(tasks) < count_tasks {
			actual_tasks, err = m.repository.GetActualLimitTasks(count_tasks - len(tasks))
			if err != nil {
				log.Debug("Failed to get tasks", sl.Err(err))
			} else {
				tasks = append(tasks, actual_tasks...)
			}
		}

		// tasks, err := m.repository.GetActualLimitTasks(count_tasks)

		// if err != nil {
		// 	log.Debug("Failed to get tasks", sl.Err(err))
		// 	return
		// }

		for _, task := range tasks {
			m.queue_tasks.Add(task)
		}

		var wg sync.WaitGroup

		maxProcs := 5
		count := maxProcs

		wg.Add(count)
		runtime.GOMAXPROCS(maxProcs)

		for i := 0; i < count; i++ {
			i := i
			go func() {
				defer wg.Done()
				log.Info(fmt.Sprintf("work #%d", i))
				m.work()
			}()
		}

		wg.Wait()
	}
}

func (m *Manager) work() {

	m.mu.Lock()
	task, task_err := m.queue_tasks.Pop()
	m.mu.Unlock()

	sites := make([]site.Site, 0)
	for task_err == nil {

		site, err := site.NewSite(task.SiteURL)

		if err != nil {
			log.Debug(fmt.Sprintf("Failed to create new site: %s", task.SiteURL), sl.Err(err))
			continue
		}

		site.Id = task.Id
		sites = append(sites, *site)

		m.mu.Lock()
		task, task_err = m.queue_tasks.Pop()
		m.mu.Unlock()
	}

	if len(sites) == 0 {
		log.Debug(fmt.Sprintf("Count sites: %d", len(sites)))
		return
	}

	robot := robot.NewRobot(m.repository)
	robot.AddList(sites)
	new_sites, success_ids, unsuccess_ids := robot.Work()

	if len(new_sites) > 0 {
		err := m.repository.AddTask(new_sites)
		if err != nil {
			log.Debug("Failed to add task", sl.Err(err))
			return
		}
	}

	if len(success_ids) > 0 {
		err := m.repository.CompleteTasks(success_ids)
		if err != nil {
			log.Debug("Failed to complete tasks", sl.Err(err))
		}
	}

	if len(unsuccess_ids) > 0 {
		err := m.repository.CompleteWithError(unsuccess_ids)
		if err != nil {
			log.Debug("Failed to point error tasks", sl.Err(err))
		}
	}

	log.Info(fmt.Sprintf("c. sites: %d c. success_ids: %d c. unsuccess_ids: %d", len(new_sites), len(success_ids), len(unsuccess_ids)))
}
