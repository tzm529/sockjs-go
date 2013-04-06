package sockjs

import (
	"net/http"
	"sync"
)

// ServeMux is SockJS-compatible HTTP request multiplexer, similar to http.ServeMux,
// but just for SockJS handlers. It can optionally wrap an alternate http.Handler which is called 
// for non-SockJS paths.
type ServeMux struct {
	mu  sync.RWMutex
	m   map[string]http.Handler
	alt http.Handler
}

// NewServeMux creates a new ServeMux with the given alternate handler.
// If alt is nil, alternate handler is not used.
func NewServeMux(alt http.Handler) *ServeMux {
	m := new(ServeMux)
	m.m = make(map[string]http.Handler)
	m.alt = alt
	return m
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := m.match(r.URL.Path)
	h.ServeHTTP(w, r)
}

func (m *ServeMux) Handle(prefix string, hfunc func(Session), c Config) {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("sockjs: prefix must not end with a slash")
	}
	if _, ok := m.m[prefix]; ok {
		panic("sockjs: multiple registrations for " + prefix)
	}

	m.mu.Lock()
	m.m[prefix] = newHandler(prefix, hfunc, &c)
	m.mu.Unlock()
}

// Does path match prefix?
func pathMatch(prefix, path string) bool {
	return len(path) >= len(prefix) && path[0:len(prefix)] == prefix
}

// Return a handler from the handler map that matches the given a path.
// Most-specific (longest) prefix wins.
// If no handler is found, return the alternate handler or http.NotFoundHandler().
func (m *ServeMux) match(path string) (h http.Handler) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var n = 0
	for k, v := range m.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v
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
