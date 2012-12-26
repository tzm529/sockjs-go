package sockjs

import (
	"net/http"
	"io"
	"encoding/json"
)

//* xhrPolling

type xhrPollingProtocol struct{}

func (p xhrPollingProtocol) contentType() string { return "application/javascript; charset=UTF-8" }

func (p xhrPollingProtocol) writeOpen(w io.Writer) (err error) {
	_, err = w.Write([]byte("o\n"))
	return
}

func (p xhrPollingProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) {
	n, err = w.Write(frame("", "\n", m...))
	return
}

func (p xhrPollingProtocol) writeClose(w io.Writer, code int, m string) {
	w.Write(cframe("", code, m, "\n"))
}

//* xhrStreaming

var prelude []byte = make([]byte, 2049)

func init() {
	for i := 0; i < 2048; i++ {
		prelude[i] = byte('h')
	}
	prelude[2048] = byte('\n')
}

type xhrStreamingProtocol struct{}

func (p xhrStreamingProtocol) contentType() string { return "application/javascript; charset=UTF-8" }

func (p xhrStreamingProtocol) writePrelude(w io.Writer) (err error) {
	_, err = w.Write(prelude)
	return
}

func (p xhrStreamingProtocol) writeOpen(w io.Writer) (err error) {
	_, err = w.Write([]byte("o\n"))
	return
}

func (p xhrStreamingProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) {
	n, err = w.Write(frame("", "\n", m...))
	return
}

func (p xhrStreamingProtocol) writeClose(w io.Writer, code int, m string) {
	w.Write(cframe("", code, m, "\n"))
}

func xhrSendHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "text/plain; charset=UTF-8")
	preflight(header, r)
	disableCache(header)

	s := h.pool.get(sessid)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	var messages []string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&messages); err != nil {
		if err == io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Payload expected."))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Broken JSON encoding."))
		return
	}

	for _, v := range messages {
		s.in.push([]byte(v))
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func xhrOptionsHandler(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Access-Control-Allow-Methods", "OPTIONS, POST")
	preflight(h, r)
	enableCache(h)
	w.WriteHeader(http.StatusNoContent)
}
