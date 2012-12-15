package sockjs

import (
	"container/list"
	"sync"
	"errors"
)

var errQueueClosed error = errors.New("queue is closed")
var errQueueWait error = errors.New("this queue forbids concurrent wait")

// Infinite message queue
type queue struct {
	*list.List  
	sync.Mutex
	*sync.Cond
	closed bool
	wait bool // forbid concurrent wait?
}

func newQueue(wait bool) (q *queue) {
	q = new(queue)
	q.List = list.New()
	q.Cond = sync.NewCond(q)
	q.wait = wait
	return
}

// Pull returns a message from the message queue or an error, if any.
// Blocks, if queue is empty.
// errQueueClosed is returned, if the queue is closed. 
// errQueueWait is returned, 
// if another goroutine is waiting and the queue does not allow concurrent waits.
func (q *queue) Pull() (m []byte, err error) {
	q.Lock()
	defer q.Unlock()
	for q.Len() == 0 { 
		if !q.closed {
			if !q.wait {
				q.Wait()
			} else {
				return nil, errQueueWait
			}
		} else {
			return nil, errQueueClosed
		}
	}
	m, _ = q.Remove(q.Front()).([]byte)
	return m, nil
}

// PullAll is like Pull except it returns all messages from the message queue.
func (q *queue) PullAll() (elems [][]byte, err error) {
	q.Lock()
	defer q.Unlock()
	for q.Len() == 0 { 
		if !q.closed { 
			if !q.wait {
				q.Wait()
			} else {
				return nil, errQueueWait
			}
		} else {
			return nil, errQueueClosed
		}
	}
	elems = make([][]byte, q.Len())
	for e, i := q.Front(), 0; e != nil; e, i = q.Front(), i+1 {
		elems[i], _ = q.Remove(e).([]byte)
	}
	return elems, nil
}

// Push pushes a message to the message queue.
// Panics, if the queue is closed.
func (q *queue) Push(m []byte) {
	q.Lock()
	defer q.Unlock()
	if q.closed { panic(errQueueClosed) }
	q.PushBack(m)
	q.Signal()
}

// Close empties the queue, marks it closed and wakes up remaining goroutines waiting on pull.
func (q *queue) Close() {
	q.Lock()
	defer q.Unlock()
	q.Init()
	q.closed = true
	q.Broadcast()
}
