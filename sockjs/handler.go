package sockjs

import (
	"net/http"
	"regexp"
)

var reInfo = regexp.MustCompile(`^/info$`)
var reIframe = regexp.MustCompile(`^/iframe[\w\d-\. ]*\.html$`)
var reSessionUrl = regexp.MustCompile(
	`^/(?:[\w- ]+)/([\w- ]+)/(xhr|xhr_send|xhr_streaming|eventsource|websocket|jsonp|jsonp_send)$`)
var reRawWebsocket = regexp.MustCompile(`^/websocket$`)

type Handler struct {
	prefix string
	hfunc  func(*Session)
	config Config
}

func newHandler(prefix string, hfunc func(*Session), c Config) (h *Handler) {
	h = new(Handler)
	h.prefix = prefix
	h.hfunc = hfunc
	h.config = c
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(h.prefix):]
	method := r.Method
	println("ServeHTTP:", path, method)

	switch {
	case method == "GET" && path == "" || path == "/":
		handleGreeting(w)
	case method == "GET" && reInfo.MatchString(path):
		handleInfo(w, r, h)
	case method == "OPTIONS" && reInfo.MatchString(path):
		handleInfoOptions(w, r)
	case method == "GET" && reIframe.MatchString(path):
		handleIframe(w, r, h)
	case method == "GET" && reSessionUrl.MatchString(path):
		matches := reSessionUrl.FindStringSubmatch(path)
		protocol := matches[2]
		switch protocol {
		case "websocket":
			handleWebsocket(w, r, h)
		}
	case method == "GET" && reRawWebsocket.MatchString(path):
		handleRawWebsocket(w, r, h)
	default:
		http.NotFound(w, r)
	}
}
