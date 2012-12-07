package sockjs

import (
	"net/http"
	"sync"
)

// ServeMux is sockjs-compatible HTTP request multiplexer, similar to http.ServeMux,
// but just for sockjs.Handlers. It can optionally wrap a http.Handler which is called for
// non-sockjs paths.
type ServeMux struct {
	mu  sync.RWMutex
	m   map[string]*Handler
	alt http.Handler
}

func NewServeMux(alt http.Handler) *ServeMux {
	m := new(ServeMux)
	m.m = make(map[string]*Handler)
	m.alt = alt
	return m
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	h = m.match(r.URL.Path)
	if h == nil {
		if m.alt != nil {
			m.alt.ServeHTTP(w, r)
		} else {
			h = http.NotFoundHandler()
		}
	} else {
		h.ServeHTTP(w, r)
	}
}

func (m *ServeMux) Handle(handler *Handler) {
	prefix := handler.config.Prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("sockjs: prefix must not end with a slash")
	}
	if _, ok := m.m[prefix]; ok {
		panic("sockjs: multiple registrations for " + prefix)
	}

	m.mu.Lock()
	m.m[prefix] = handler
	m.mu.Unlock()
}

// Does path match prefix?
func pathMatch(prefix, path string) bool {
	return len(path) >= len(prefix) && path[0:len(prefix)] == prefix
}

// Find a handler on a handler map given a path string
// Most-specific (longest) prefix wins
func (m *ServeMux) match(path string) (h http.Handler) {
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
	return
}
