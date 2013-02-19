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

func noCache(h http.Header) {
	h.Add("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
}

func xhrCors(h http.Header, r *http.Request) {
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

// Some load balancers do sticky sessions, but only if there is
// a JSESSIONID cookie. If this cookie isn't yet set, we shall
// set it to a dummy value. It doesn't really matter what, as
// session information is usually added by the load balancer.
func sid(h *handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Jsessionid {
		return
	}

	// Users can supply a function
	if h.config.JsessionidFunc != nil {
		h.config.JsessionidFunc(w, r)
		return
	}

	// We need to set it every time, to give the loadbalancer
	// opportunity to attach its own cookies.
	jsid, err := r.Cookie(dummyCookie.Name)
	if err == http.ErrNoCookie {
		jsid = dummyCookie
	}
	jsid.Path = "/"
	http.SetCookie(w, jsid)
}
