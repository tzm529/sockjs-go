// chat room example

package main

import (
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"net/http"
	"strings"
	"sync"
)

// ChatPool is a structure for holding chat users and broadcasting messages to them.
type userPool struct {
	sync.Mutex
	users map[sockjs.Session]struct{}
}

func newUserPool() (p *userPool) {
	p = new(userPool)
	p.users = make(map[sockjs.Session]struct{})
	return
}

func (p *userPool) add(s sockjs.Session) {
	p.Lock()
	defer p.Unlock()
	p.users[s] = struct{}{}
}

func (p *userPool) remove(s sockjs.Session) {
	p.Lock()
	defer p.Unlock()
	delete(p.users, s)
}

func (p *userPool) broadcast(m []byte) {
	p.Lock()
	defer p.Unlock()
	for s := range p.users {
		s.Send(m)
	}
}

var users *userPool = newUserPool()

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
