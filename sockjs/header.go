package sockjs

//* Tools for manipulating HTTP headers.

import (
	"fmt"
	"net/http"
	"time"
)

func addExpires(h http.Header) {
	h.Add("Cache-Control", fmt.Sprintf("public, max-age=%d", 365*24*60*60))
	h.Add("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))
	h.Add("Access-Control-Max-Age", fmt.Sprintf("%d", 365*24*60*60))
}
