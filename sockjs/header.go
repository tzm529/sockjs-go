package sockjs

//* Tools for manipulating HTTP headers.

import (
	"fmt"
	"net/http"
	"time"
)

func expires(h http.Header) {
	h.Add("Cache-Control", fmt.Sprintf("public, max-age=%d", 365*24*60*60))
	h.Add("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))
	h.Add("Access-Control-Max-Age", fmt.Sprintf("%d", 365*24*60*60))
}

func cors(h http.Header, r *http.Request) {
	h.Add("Access-Control-Allow-Credentials", "true")
	h.Add("Access-Control-Allow-Origin", getOriginHeader(r))
	allowHeaders := r.Header.Get("Access-Control-Request-Headers")
	if allowHeaders != "" && allowHeaders != "null" {
		h.Add("Access-Control-Allow-Headers", allowHeaders)
	}
}

func corsAllowMethods(h http.Header, r *http.Request, methods string) {
	cors(h, r)
	h.Add("Access-Control-Allow-Methods", methods)
}

func disableCache(h http.Header) {
	h.Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
}

func getOriginHeader(r *http.Request) string {
	origin := r.Header.Get("Origin")
	if origin == "" || origin == "null" {
		origin = "*"
	}
	return origin
}
