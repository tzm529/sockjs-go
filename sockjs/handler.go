package sockjs

import (
	"net/http"
)

type Handler struct {
	config Config
}

func NewHandler(c Config) (h *Handler) {
	h = new(Handler)
	return h
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	println("ServeHTTP:", path, method)

	switch {
	default:
		http.NotFound(w, r)
	}

}
