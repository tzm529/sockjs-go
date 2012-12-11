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
	pool *sessionPool
}

func newHandler(prefix string, hfunc func(*Session), c Config) (h *Handler) {
	h = new(Handler)
	h.prefix = prefix
	h.hfunc = hfunc
	h.config = c
	h.pool = newSessionPool()
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
		handleInfo(h, w, r)
	case method == "OPTIONS" && reInfo.MatchString(path):
		handleInfoOptions(w, r)
	case method == "GET" && reIframe.MatchString(path):
		handleIframe(h, w, r)
	case method == "GET" && reRawWebsocket.MatchString(path):
		handleRawWebsocket(h, w, r)
	case method == "GET" && reSessionUrl.MatchString(path):
		matches := reSessionUrl.FindStringSubmatch(path)
		protocol := matches[2]
		switch protocol {
		case "websocket":
			handleWebsocket(h, w, r)
		}
	case method == "POST" && reSessionUrl.MatchString(path):
		matches := reSessionUrl.FindStringSubmatch(path)
		//sessid := matches[1]
		protocol := matches[2]
		switch protocol {
		case "websocket":
			handleWebsocketPost(w, r)
		case "xhr":
			//handleXhrPolling(w, r, sessid, h)
		case "xhr_send":
			//handleXhrSend(w, r, sessid, h)
		case "xhr_streaming":
			//xhrStreamingHandler(w, r, sessid, h)
		case "jsonp_send":
			//handleJsonpSend(w, r, sessid, h)
		}
	case method == "OPTIONS" && reSessionUrl.MatchString(path):
		//xhrHandlerOptions(w, r)
	default:
		http.NotFound(w, r)
	}
}
