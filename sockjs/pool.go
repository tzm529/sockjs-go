package sockjs

import (
	"sync"
)

// Session pool for legacy protocols.
type pool struct {
	sync.RWMutex
	pool map[string]*legacySession
}

func newPool() (p *pool) {
	p = new(pool)
	p.pool = make(map[string]*legacySession)
	return
}

func (p *pool) get(sessid string) *legacySession {
	p.RLock()
	defer p.RUnlock()
	return p.pool[sessid]
}

func (p *pool) getOrCreate(sessid string) (s *legacySession, exists bool) {
	p.Lock()
	defer p.Unlock()
	s, exists = p.pool[sessid]
	if exists {
		return
	}
	p.pool[sessid] = new(legacySession)
	s = p.pool[sessid]
	return
}

func (p *pool) remove(sessid string) (s *legacySession) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}
