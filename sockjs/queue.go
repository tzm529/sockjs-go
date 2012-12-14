package sockjs

import (
	"container/list"
	"sync"
	"errors"
)

var errQueueClosed error = errors.New("queue is closed")

// Infinite message queue
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

// Pull returns a message from the message queue or nil, if the queue is closed. 
// Blocks, if queue is empty.
func (q *queue) pull() []byte {
	q.Lock()
	defer q.Unlock()
	for q.Len() == 0 { 
		if !q.closed { 
			q.Wait()
		} else {
			return nil
		}
	}
	m, _ := q.Remove(q.Front()).([]byte)
	return m
}

// PullAll returns all messages from the message queue or nil, if the queue is closed. 
// Blocks, if queue is empty.
func (q *queue) pullAll() [][]byte {
	q.Lock()
	defer q.Unlock()
	for q.Len() == 0 { 
		if !q.closed { 
			q.Wait()
		} else {
			return nil
		}
	}
	elems := make([][]byte, q.Len())
	for e, i := q.Front(), 0; e != nil; e, i = q.Front(), i+1 {
		elems[i], _ = q.Remove(e).([]byte)
	}
	return elems
}

// Push pushes a message to the message queue.
// Panics, if the queue is closed.
func (q *queue) push(m []byte) {
	q.Lock()
	defer q.Unlock()
	if q.closed { panic(errQueueClosed) }
	q.PushBack(m)
	q.Signal()
}

// Close empties the queue, marks it closed and wakes up remaining goroutines waiting on pull.
func (q *queue) close() {
	q.Lock()
	defer q.Unlock()
	q.Init()
	q.closed = true
	q.Broadcast()
}
