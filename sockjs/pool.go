package sockjs

import (
	"sync"
)

// SessionPool is a structure for thread-safely storing sessions and broadcasting messages to them.
type SessionPool struct {
	mu   sync.RWMutex
	pool map[Session]struct{}
}

func NewSessionPool() (p *SessionPool) {
	p = new(SessionPool)
	p.pool = make(map[Session]struct{})
	return
}

// Add adds the given session to the session pool.
func (p *SessionPool) Add(s Session) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool[s] = struct{}{}
}

// Remove removes the given session from the session pool.
// It is safe to remove non-existing sessions.
func (p *SessionPool) Remove(s Session) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pool, s)
}

// Broadcast sends the given message to every session in the pool.
func (p *SessionPool) Broadcast(m []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for s := range p.pool {
		s.Send(m)
	}
}

// Session pool for storing legacy sessions.
type legacyPool struct {
	sync.RWMutex
	pool map[string]*legacySession
}

func newLegacyPool() (p *legacyPool) {
	p = new(legacyPool)
	p.pool = make(map[string]*legacySession)
	return
}

func (p *legacyPool) get(sessid string) *legacySession {
	p.RLock()
	defer p.RUnlock()
	return p.pool[sessid]
}

func (p *legacyPool) getOrCreate(sessid string) (s *legacySession, exists bool) {
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

func (p *legacyPool) remove(sessid string) (s *legacySession) {
	p.Lock()
	defer p.Unlock()
	s = p.pool[sessid]
	delete(p.pool, sessid)
	return
}
