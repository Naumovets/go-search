package robot

// Пакет поискового робота, обрабатывает очереди из Site

import (
	"sync"

	"github.com/Naumovets/go-search/internal/collections/queue"
	log "github.com/Naumovets/go-search/internal/logger"
	"github.com/Naumovets/go-search/internal/logger/sl"
	"github.com/Naumovets/go-search/internal/repositories/tasks"
	"github.com/Naumovets/go-search/internal/site"
)

type Robot struct {
	queue        queue.Queue[site.Site]
	newSites     []site.Site
	successIds   []int
	unsuccessIds []int
	mu           sync.Mutex
	repository   *tasks.Repository
}

func NewRobot(rep *tasks.Repository) *Robot {
	return &Robot{
		repository: rep,
		queue:      *queue.NewQueue[site.Site](),
		successIds: make([]int, 0),
	}
}

func (r *Robot) AddOne(s *site.Site) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queue.Add(s)
}

func (r *Robot) AddList(l []*site.Site) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range l {
		r.queue.Add(item)
	}
}

func (r *Robot) Work() ([]site.Site, []int, []int) {

	for r.queue.Count != 0 {

		site, err := r.queue.Pop()
		if err != nil {
			log.Debug("Failed to get task", sl.Err(err))
			continue
		}

		newSites, err := site.Analys()

		if err != nil {
			log.Debug("Failed to analys site", sl.Err(err))
			r.unsuccessIds = append(r.unsuccessIds, site.Id)
			continue
		}
		r.successIds = append(r.successIds, site.Id)
		r.newSites = append(r.newSites, newSites...)
	}
	return r.newSites, r.successIds, r.unsuccessIds
}
