package sockjs

import (
	"net/http"
)

func greetingHandler(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte("Welcome to SockJS!\n"))
}
