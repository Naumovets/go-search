package robot

import (
	"github.com/Naumovets/go-search/internal/collections/queue"
	"github.com/Naumovets/go-search/internal/site"
)

type Robot struct {
	queue queue.Queue[site.Site]
}

func NewRobot() *Robot {
	return &Robot{
		queue: *queue.NewQueue[site.Site](),
	}
}

func (r *Robot) AddOne(s site.Site) {
	r.queue.Add(s)
}

func (r *Robot) AddList(l []site.Site) {
	for _, item := range l {
		r.queue.Add(item)
	}
}
