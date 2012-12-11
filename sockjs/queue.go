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
	closed bool
}

func newQueue() (q *queue) {
	q = new(queue)
	q.List = list.New()
	q.Cond = sync.NewCond(q)
	return
}

// Returns a message from the message queue or an error,
// if queue is closed. Blocks, if queue is empty.
func (q *queue) pull() ([]byte, error) {
	q.Lock()
	defer q.Unlock()

	if q.closed { return nil, ErrSessionClosed }

	if q.Len() == 0 {
		q.Wait()
	}

	if q.closed { return nil, ErrSessionClosed }
	m, _ := q.Remove(q.Front()).([]byte)
	return m, nil
}

// Pushes a message to the message queue or returns an error,
// if queue is closed.
func (q *queue) push(m []byte) error {
	q.Lock()
	defer q.Unlock()

	if q.closed { return ErrSessionClosed }
	q.PushBack(m)
	q.Signal()
	return nil
}

// Closes the message queue by waking up every goroutine blocking on pull().
func (q *queue) close() {
	q.Lock()
	defer q.Unlock()

	q.Broadcast()
	q.closed = true
}
