package sockjs

import (
	"sync"
)

type sessionFactory func() Session

// session pool
type pool struct {
	sync.RWMutex
	pool map[string]Session
}

func newPool() *pool {
	pool := new(pool)
	pool.pool = make(map[string]Session)
	return pool
}

func (p *pool) get(sessid string) Session {
	p.RLock()
	defer p.RUnlock()
	return p.pool[sessid]
}

func (p *pool) getOrCreate(sessid string, f sessionFactory) (s Session, exists bool) {
	p.Lock()
	defer p.Unlock()
	s, exists = p.pool[sessid]
	if exists { return }
	p.pool[sessid] = f()
	s = p.pool[sessid]
	return
}

func (p *pool) remove(sessid string) (s Session) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}

func (p *pool) close() {
}