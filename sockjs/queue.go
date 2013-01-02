package sockjs

import (
	"container/list"
	"errors"
	"sync"
)

var errQueueClosed error = errors.New("queue closed")

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

// Pull returns a message from the message queue or an error, if any.
// Blocks, if queue is empty.
// errQueueClosed is returned, if the queue is closed. 
func (q *queue) pull() (m []byte, err error) {
	q.Lock()
	defer q.Unlock()

	for q.Len() == 0 {
		if !q.closed {
			q.Wait()
		} else {
			return nil, errQueueClosed
		}
	}
	m, _ = q.Remove(q.Front()).([]byte)
	return m, nil
}

// PullAll is like Pull except it returns all messages from the message queue.
func (q *queue) pullAll() (messages [][]byte, err error) {
	q.Lock()
	defer q.Unlock()

	for q.Len() == 0 {
		if !q.closed {
			q.Wait()
		} else {
			return nil, errQueueClosed
		}
	}
	messages = make([][]byte, q.Len())
	for e, i := q.Front(), 0; e != nil; e, i = q.Front(), i+1 {
		messages[i], _ = q.Remove(e).([]byte)
	}
	return messages, nil
}

// PullNow is like Pull except it does not block.
func (q *queue) pullNow() (m []byte, err error) {
	q.Lock()
	defer q.Unlock()

	if q.closed {
		return nil, errQueueClosed
	}
	if q.Len() == 0 {
		return
	}
	m, _ = q.Remove(q.Front()).([]byte)
	return m, nil
}

// Push pushes given messages to the message queue.
// Panics, if the queue is closed.
func (q *queue) push(messages ...[]byte) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		panic(errQueueClosed)
	}
	for _, v := range messages {
		q.PushBack(v)
	}
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
