package common

import (
	"testing"
)

func Test_Queue(t *testing.T) {
	queue := Queue{}

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	t.Log("Queue size:", queue.Size())
	t.Log("Front element:", queue.Front())

	t.Log("Dequeued element:", queue.Dequeue())

	t.Log("Queue size:", queue.Size())
	t.Log("Front element:", queue.Front())

	t.Log("Is queue empty?", queue.IsEmpty())

	structQueue := Queue{}
	structQueue.Enqueue(PoolElem{
		ReqID:     "a",
		RandomNum: 0,
	})
	structQueue.Enqueue(PoolElem{
		ReqID:     "a",
		RandomNum: 1,
	})
	structQueue.Enqueue(PoolElem{
		ReqID:     "b",
		RandomNum: 3,
	})
	t.Log("Queue size:", structQueue.Size())
	t.Log("Front element:", structQueue.Front())

	first := structQueue.Dequeue()
	t.Log("Dequeued element: first", first)
	t.Log("Dequeued element: first reqID", structQueue.Dequeue().(PoolElem).ReqID)

	t.Log("Dequeued element:", structQueue.Dequeue())
	t.Log("Dequeued element:", structQueue.Dequeue())
}

func Test_QueueMany(t *testing.T) {
	q := Queue{}
	q.Enqueue(PoolElem{
		ReqID:     "a",
		RandomNum: 1,
	})
	q.Enqueue(PoolElem{
		ReqID:     "a",
		RandomNum: 2,
	})
	q.Enqueue(PoolElem{
		ReqID:     "b",
		RandomNum: 3,
	})
	q.Enqueue(PoolElem{
		ReqID:     "c",
		RandomNum: 4,
	})

	e1 := q.Dequeue().(PoolElem)
	t.Logf("%+v", e1)

	es := q.DequeueMulti(4)
	t.Logf("%+v", es)

	es = q.DequeueMulti(3)
	t.Logf("%+v", es)

	t.Logf("queue len %d", q.Size())
}
