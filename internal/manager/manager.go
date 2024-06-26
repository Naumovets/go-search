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
	"github.com/Naumovets/go-search/internal/repositories/db"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	"github.com/Naumovets/go-search/internal/robot"
	"github.com/Naumovets/go-search/internal/site"
)

type Manager struct {
	queueTasks *queue.Queue[entities.Task]
	taskRep    *tasks.Repository
	dbRep      *db.Repository
}

func NewManager(taskRep *tasks.Repository, dbRep *db.Repository) *Manager {
	return &Manager{
		queueTasks: queue.NewQueue[entities.Task](),
		taskRep:    taskRep,
		dbRep:      dbRep,
	}
}

func (m *Manager) Start(countTasks int) {

	for {
		// get tasks to queue
		tasks, err := m.taskRep.GetNoCompleteLimitTasks(countTasks)
		var actualTasks []entities.Task

		if err != nil {
			log.Debug("Failed to get no complete tasks", sl.Err(err))
			actualTasks, err = m.taskRep.GetActualLimitTasks(countTasks)
			if err != nil {
				log.Debug("Failed to get tasks", sl.Err(err))
				return
			}
			tasks = actualTasks

		} else if len(tasks) < countTasks {
			actualTasks, err = m.taskRep.GetActualLimitTasks(countTasks - len(tasks))
			if err != nil {
				log.Debug("Failed to get tasks", sl.Err(err))
			} else {
				tasks = append(tasks, actualTasks...)
			}
		}

		// tasks, err := m.taskRep.GetActualLimitTasks(countTasks)

		// if err != nil {
		// 	log.Debug("Failed to get tasks", sl.Err(err))
		// 	return
		// }

		for _, task := range tasks {
			m.queueTasks.Add(&task)
		}

		var wg sync.WaitGroup

		maxProcs := 7
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

	m.queueTasks.Mu.Lock()
	task, taskErr := m.queueTasks.Pop()
	m.queueTasks.Mu.Unlock()

	sites := make([]*site.Site, 0)
	for taskErr == nil {

		site, err := site.NewSite(task.SiteURL)

		if err != nil {
			log.Debug(fmt.Sprintf("Failed to create new site: %s", task.SiteURL), sl.Err(err))
			continue
		}

		site.Id = task.Id
		sites = append(sites, site)

		m.queueTasks.Mu.Lock()
		task, taskErr = m.queueTasks.Pop()
		m.queueTasks.Mu.Unlock()
	}

	if len(sites) == 0 {
		log.Debug(fmt.Sprintf("Count sites: %d", len(sites)))
		return
	}

	robot := robot.NewRobot(m.taskRep)
	robot.AddList(sites)
	newSites, successIds, unsuccessIds := robot.Work()

	if len(newSites) > 0 {
		err := m.taskRep.AddTask(newSites)
		if err != nil {
			log.Debug("Failed to add task", sl.Err(err))
			return
		}
	}

	if len(successIds) > 0 {
		successSites := getCompletedSites(sites, successIds)
		err := m.dbRep.AddSites(successSites)

		if err != nil {
			log.Debug("Failed to add sites to db", sl.Err(err))
		}

		err = m.taskRep.CompleteTasks(successIds)
		if err != nil {
			log.Debug("Failed to complete tasks", sl.Err(err))
		}
	}

	if len(unsuccessIds) > 0 {
		err := m.taskRep.CompleteWithError(unsuccessIds)
		if err != nil {
			log.Debug("Failed to point error tasks", sl.Err(err))
		}
	}

	log.Info(fmt.Sprintf("c. sites: %d c. successIds: %d c. unsuccessIds: %d", len(newSites), len(successIds), len(unsuccessIds)))
}
