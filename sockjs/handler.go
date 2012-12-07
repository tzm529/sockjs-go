package sockjs

import (
	"net/http"
	"regexp"
)

var reSessionUrl = regexp.MustCompile(
	`/(?:[\w- ]+)/([\w- ]+)/(xhr|xhr_send|xhr_streaming|eventsource|websocket|jsonp|jsonp_send)$`)

type Handler struct {
	OnOpen func(*Conn)
	OnMessage func(*Conn, string)
	OnClose func(*Conn)
	config Config
}

func NewHandler(c Config) (h *Handler) {
	h = new(Handler)
	h.config = c
	return h
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	println("ServeHTTP:", path, method)

	switch {
	case method == "GET" && reSessionUrl.MatchString(path):
		matches := reSessionUrl.FindStringSubmatch(path)
		protocol := matches[2]
		switch protocol {
		case "websocket":
			handleWebsocket(w, r, s)
		}
	default:
		http.NotFound(w, r)
	}

}
