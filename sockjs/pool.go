package sockjs

import (
	"time"
	"sync"
)

type poolElem struct {
	sync.Mutex
	*sync.Cond
	s *session
	onloan bool
}

//* pool for non-websocket sessions
type pool struct {
	sync.RWMutex
	pool map[string]*poolElem
}

func newPool() (p *pool) {
	p := new(p)
	p.p = make(map[string]*poolElem)
	return
}

func (p *pool) restore(sessid string) {
	p.RLock()
	pe, exists := p.pool[sessid]
	p.RUnlock()

	if !exists { return }
	pe.Lock()
	defer pe.Unlock()
	pe.onloan = false
	pe.Signal()
}

func (p *pool) get(sessid string) *session {
	p.RLock()
	pe := p.pool[sessid]
	p.RUnlock()

	pe.Lock()
	defer pe.Unlock()
	for !pe.onloan {
		pe.Wait()
	}

	pe.onloan = true
	return pe.s
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
