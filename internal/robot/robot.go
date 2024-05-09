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
	queue       queue.Queue[site.Site]
	new_sites   []site.Site
	success_ids []int
	mu          sync.Mutex
	repository  *tasks.Repository
}

func NewRobot(rep *tasks.Repository) *Robot {
	return &Robot{
		repository:  rep,
		queue:       *queue.NewQueue[site.Site](),
		success_ids: make([]int, 0),
	}
}

func (r *Robot) AddOne(s site.Site) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queue.Add(s)
}

func (r *Robot) AddList(l []site.Site) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range l {
		r.queue.Add(item)
	}
}

func (r *Robot) Work() ([]site.Site, []int) {

	for r.queue.Count != 0 {

		site, err := r.queue.Pop()
		if err != nil {
			log.Debug("Failed to get task", sl.Err(err))
			continue
		}

		new_sites, err := site.Analys()

		if err != nil {
			log.Debug("Failed to analys site", sl.Err(err))
			continue
		}
		r.success_ids = append(r.success_ids, site.Id)
		r.new_sites = append(r.new_sites, new_sites...)
	}
	return r.new_sites, r.success_ids
}
