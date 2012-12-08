package sockjs

import (
	"net/http"
)

func handleGreeting(w http.ResponseWriter) {
	h := w.Header()
	addExpires(h)
	h.Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte("Welcome to SockJS!\n"))
}
