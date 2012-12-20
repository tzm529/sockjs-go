package sockjs

//* Tools for manipulating HTTP headers.

import (
	"fmt"
	"net/http"
	"time"
)

func enableCache(h http.Header) {
	h.Add("Cache-Control", fmt.Sprintf("public, max-age=%d", 365*24*60*60))
	h.Add("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))
	h.Add("Access-Control-Max-Age", fmt.Sprintf("%d", 365*24*60*60))
}

func disableCache(h http.Header) {
	h.Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
}

func preflight(h http.Header, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "null" {
		origin = "*"
	}
	h.Add("Access-Control-Allow-Origin", origin)
	headers := r.Header.Get("Access-Control-Request-Headers")
	if headers != "" {
		h.Add("Access-Control-Allow-Headers", headers)
	}

	h.Add("Access-Control-Allow-Credentials", "true")
}
