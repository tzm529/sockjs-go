package sockjs

import (
	"net/url"
	"net/http"
)

type Protocol uint8

const(
	ProtocolRawWebsocket Protocol = iota
	ProtocolWebsocket
	ProtocolXhrPolling
	ProtocolXhrStreaming
	ProtocolJsonp
	ProtocolEventSource
	ProtocolHtmlfile
)

// RequestInfo contains information copied from the last received request associated with the session.
type RequestInfo struct {
	// URL which was sought.
    URL    url.URL

	// Copy of the HTTP headers listed in Config.Headers.
    Header http.Header

	// Host on which the URL was sought.
    Host string

	// Remote address of the client in "IP:port" format.
    RemoteAddr string

	// RequestURI is the unmodified Request-URI of the
    // Request-Line (RFC 2616, Section 5.1) as sent by the client
    // to a server. Usually the URL field should be used instead.
    RequestURI string

	// Prefix of the URL on which the request was handled.
	Prefix string
}

func newRequestInfo(r *http.Request, 
	prefix string, 
	headers []string) (info *RequestInfo) {
	info = new(RequestInfo)
	info.URL = *r.URL
	info.Host = r.Host
	info.RequestURI = r.RequestURI
	info.RemoteAddr = r.RemoteAddr
	info.Prefix = prefix

	h := r.Header
	for _,k := range headers {
		k = http.CanonicalHeaderKey(k)
		klen := len(h[k])
		if klen > 0 {
			info.Header[k] = make([]string, klen)
			copy(info.Header[k], h[k])
		}
	}
	return
}
