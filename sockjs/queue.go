package sockjs

import (
	"container/list"
	"sync"
)

// Message queue
type queue struct {
	*list.List  
	sync.Mutex
	*sync.Cond
}

func newQueue() (q *queue) {
	q = new(queue)
	q.List = list.New()
	q.Cond = sync.NewCond(q)
	return
}

// Pulls a message from the message queue.
// Blocks, if queue is empty.
func (q *queue) pull() []byte {
	q.Lock()
	defer q.Unlock()

	if q.Len() == 0 {
		q.Wait()
	}
	e := q.Front()
	if e == nil {
		return nil
	}
	m, _ := q.Remove(e).([]byte)
	return m
}

// Pushes a message to the message queue.
func (q *queue) push(m []byte) {
	q.Lock()
	defer q.Unlock()
	q.PushBack(m)
	q.Signal()
}
