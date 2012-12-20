package sockjs

import (
	"net/http"
	"sync"
)

// Server is sockjs-compatible HTTP request multiplexer, similar to http.ServeMux,
// but just for sockjs.Handlers. It can optionally wrap an alternate http.Handler which is called 
// for non-sockjs paths.
type Server struct {
	mu   sync.RWMutex
	m    map[string]http.Handler
	alt  http.Handler
	pool *pool
}

func NewServer(alt http.Handler) *Server {
	m := new(Server)
	m.m = make(map[string]http.Handler)
	m.alt = alt
	m.pool = newPool()
	return m
}

func (m *Server) Close() {
	m.pool.close()
}

func (m *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := m.match(r.URL.Path)
	h.ServeHTTP(w, r)
}

func (m *Server) Handle(prefix string, hfunc func(Session), c Config) {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("sockjs: prefix must not end with a slash")
	}
	if _, ok := m.m[prefix]; ok {
		panic("sockjs: multiple registrations for " + prefix)
	}

	m.mu.Lock()
	m.m[prefix] = newHandler(m.pool, prefix, hfunc, c)
	m.mu.Unlock()
}

// Does path match prefix?
func pathMatch(prefix, path string) bool {
	return len(path) >= len(prefix) && path[0:len(prefix)] == prefix
}

// Return a handler from the handler map that matches the given a path.
// Most-specific (longest) prefix wins.
// If no handler is found, return the alternate handler or http.NotFoundHandler().
func (m *Server) match(path string) (h http.Handler) {
	var n = 0
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = http.Handler(v)
		}
	}
	if h == nil {
		if m.alt != nil {
			return m.alt
		} else {
			h = http.NotFoundHandler()
		}
	}
	return
}
