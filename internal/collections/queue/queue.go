package queue

import "fmt"

type Queue[T any] struct {
	first *node[T]
	last  *node[T]
	Count int
}

type node[T any] struct {
	next  *node[T]
	value T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Add(val T) {
	n := &node[T]{
		value: val,
	}
	if q.Count == 0 {
		q.first = n
		q.last = n
	}
	q.last.next = n
	q.last = n
	q.Count++
}

func (q *Queue[T]) Pop() (T, error) {
	if q.Count == 0 {
		var zero T
		return zero, fmt.Errorf("err: queue if empty")
	}

	val := q.first.value
	q.first = q.first.next
	q.Count--

	if q.Count == 0 {
		q.first = nil
		q.last = nil
	}
	return val, nil
}
