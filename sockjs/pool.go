package sockjs

import (
	"sync"
)

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
	p.pool[sessid] = s
	p.Unlock()
}

func (p *pool) remove(sessid string) (s Session) {
	p.Lock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	p.Unlock()
	return
}
