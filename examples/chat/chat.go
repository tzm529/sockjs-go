// chat room example

package main

import (
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"net/http"
	"strings"
	"sync"
)

// Pool is a structure for holding chat users and broadcasting messages to them.
type pool struct {
	sync.RWMutex
	pool map[sockjs.Session]struct{}
}

func newPool() (p *pool) {
	p = new(pool)
	p.pool = make(map[sockjs.Session]struct{})
	return
}

func (p *pool) add(s sockjs.Session) {
	p.Lock()
	defer p.Unlock()
	p.pool[s] = struct{}{}
}

func (p *pool) remove(s sockjs.Session) {
	p.Lock()
	defer p.Unlock()
	delete(p.pool, s)
}

func (p *pool) broadcast(m []byte) {
	p.RLock()
	defer p.RUnlock()
	for s := range p.pool {
		s.Send(m)
	}
}

var users *pool = newPool()

func chatHandler(s sockjs.Session) {
	users.add(s)
	defer users.remove(s)

	for {
		m := s.Receive()
		if m == nil {
			break
		}
		fullAddr := s.Info().RemoteAddr
		addr := fullAddr[:strings.LastIndex(fullAddr, ":")]
		m = []byte(fmt.Sprintf("%s: %s", addr, m))
		users.broadcast(m)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	server := sockjs.NewServer(http.DefaultServeMux)
	conf := sockjs.NewConfig()
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/", indexHandler)
	server.Handle("/chat", chatHandler, conf)

	err := http.ListenAndServe(":8081", server)
	if err != nil {
		fmt.Println(err)
	}
}
