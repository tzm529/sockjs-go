package sockjs

import (
	"sync"
)

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

func (p *pool) set(sessid string, s Session) {
	p.Lock()
	defer p.Unlock()
	p.pool[sessid] = s
}

func (p *pool) remove(sessid string) (s Session) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}
