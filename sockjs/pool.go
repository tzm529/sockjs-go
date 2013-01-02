package sockjs

import (
	"time"
	"sync"
)

//* pool for "Protocol"/non-websocket sessions
type pool struct {
	sync.RWMutex
	disconnectDelay time.Duration
	pool map[string]*session
	closed bool
}

func newPool(disconnectDelay time.Duration) *pool {
	pool := new(pool)
	pool.disconnectDelay = disconnectDelay
	pool.pool = make(map[string]*session)

	go pool.garbageCollector()
	return pool
}

func (p *pool) get(sessid string) *session {
	p.RLock()
	defer p.RUnlock()
	return p.pool[sessid]
}

func (p *pool) getOrCreate(sessid string) (s *session, exists bool) {
	p.Lock()
	defer p.Unlock()
	s, exists = p.pool[sessid]
	if exists {
		return
	}
	p.pool[sessid] = new(session)
	s = p.pool[sessid]
	return
}

func (p *pool) remove(sessid string) (s *session) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}

func (p *pool) close() {
	p.Lock()
	defer p.Unlock()
	p.closed = true
}

// Garbage collector cleans up timeouted connections.
func (p *pool) garbageCollector() {
	for {
		time.Sleep(p.disconnectDelay)
		p.Lock()
		for k, v := range p.pool {
			timeouted := time.Since(v.lastRecvTime()) > p.disconnectDelay
		    if timeouted || v.closed() { 
				if timeouted { v.timeout() }
				v.cleanup()
				delete(p.pool, k)
			}
		}
		if p.closed { return }
		p.Unlock()
	}
}