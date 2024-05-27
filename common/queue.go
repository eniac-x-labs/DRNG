package common

import "sync"

// Queue first in first out
type Queue struct {
	sync.RWMutex
	items []interface{}
}

// Enqueue append elems to the queue tail
func (q *Queue) Enqueue(item ...interface{}) {
	q.items = append(q.items, item...)
}

// Dequeue get the first elem on queue head, and delete it
func (q *Queue) Dequeue() interface{} {
	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// DequeueMulti get the specified number of elements starting from head
func (q *Queue) DequeueMulti(n int) []interface{} {
	if len(q.items) < n {
		return nil
	}
	items := q.items[:n]
	q.items = q.items[n:]
	return items
}

// Front get but not delete the first element of the queue
func (q *Queue) Front() interface{} {
	if len(q.items) == 0 {
		return nil
	}
	return q.items[0]
}

// IsEmpty check whether the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size return size of queue
func (q *Queue) Size() int {
	return len(q.items)
}
