package sockjs

import (
	"sync"
)

//* interface and pool for "Protocol"/non-websocket sessions

type session interface {
	Session
	in() *queue
	out() *queue
	closed() bool
}

type sessionFactory func(pool *pool) session

type pool struct {
	sync.RWMutex
	pool map[string]session
}

func newPool() *pool {
	pool := new(pool)
	pool.pool = make(map[string]session)
	return pool
}

func (p *pool) get(sessid string) session {
	p.RLock()
	defer p.RUnlock()
	return p.pool[sessid]
}

func (p *pool) getOrCreate(sessid string, f sessionFactory) (s session, exists bool) {
	p.Lock()
	defer p.Unlock()
	s, exists = p.pool[sessid]
	if exists {
		return
	}
	p.pool[sessid] = f(p)
	s = p.pool[sessid]
	return
}

func (p *pool) remove(sessid string) (s session) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}

func (p *pool) close() {
}
