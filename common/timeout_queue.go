package common

import "time"

type TimeoutQueue struct {
	Queue Queue
}

type TimeoutInfo struct {
	Seq     int64
	Timeout time.Time
}

func (q *TimeoutQueue) Enqueue(t TimeoutInfo) {
	q.Queue.Lock()
	defer q.Queue.Unlock()
	q.Queue.Enqueue(t)
}

func (q *TimeoutQueue) Dequeue() TimeoutInfo {
	q.Queue.Lock()
	defer q.Queue.Unlock()
	inter := q.Queue.Dequeue()
	if inter == nil {
		return TimeoutInfo{}
	}
	return inter.(TimeoutInfo)
}

func (q *TimeoutQueue) FrontTime() *time.Time {
	q.Queue.RLock()
	defer q.Queue.RUnlock()
	inter := q.Queue.Front()
	if inter == nil {
		return nil
	}
	timeout := inter.(TimeoutInfo).Timeout

	return &timeout
}
