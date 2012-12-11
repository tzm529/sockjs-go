package sockjs

import (
	"container/list"
	"sync"
)

// message queue
type queue struct {
	*list.List  
	sync.Mutex
}

func newQueue() (q *queue) {
	q = new(queue)
	q.List = list.New()
	return
}

// pull a message from the message queue
func (q *queue) pull() []byte {
	q.Lock()
	defer q.Unlock()
	e := q.Front()
	if e == nil {
		return nil
	}
	m, _ := q.Remove(e).([]byte)
	return m
}

// push a message to the message queue
func (q *queue) push(m []byte) {
	q.Lock()
	defer q.Unlock()
	q.PushBack(m)
}
