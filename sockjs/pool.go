package sockjs

import (
	"sync"
)

//* pool for "Protocol"/non-websocket sessions
type pool struct {
	sync.RWMutex
	pool map[string]*session
}

func newPool() *pool {
	pool := new(pool)
	pool.pool = make(map[string]*session)
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
}
