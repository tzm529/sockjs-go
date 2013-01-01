package sockjs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type jsonpProtocol struct {
	callback string
}

func (p *jsonpProtocol) contentType() string { return "application/javascript; charset=UTF-8" }

func (p *jsonpProtocol) writeOpen(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%s(\"o\");\r\n", p.callback)
	return
}

func (p *jsonpProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) {
	js, _ := json.Marshal(string(aframe("", "", m...)))
	n, err = fmt.Fprintf(w, "%s(%s);\r\n", p.callback, js)
	return
}

func (p *jsonpProtocol) writeClose(w io.Writer, code int, m string) {
	fmt.Fprintf(w, "%s(\"c[%d,\\\"%s\\\"]\");\r\n", p.callback, code, m)
}

func (p *jsonpProtocol) protocol() Protocol       { return ProtocolJsonp }
func (p *jsonpProtocol) streaming() preludeWriter { return nil }

func jsonpHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	callback := r.Form.Get("c")
	if callback == "" {
		http.Error(w, `"callback" parameter required`, http.StatusInternalServerError)
		return
	}
	if reCallback.MatchString(callback) {
		http.Error(w, `invalid "callback" parameter`, http.StatusInternalServerError)
		return
	}

	p := new(jsonpProtocol)
	p.callback = callback
	protocolHandler(h, w, r, sessid, p)
}

func jsonpSendHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	header := w.Header()
	header.Add("Content-Type", "text/plain; charset=UTF-8")
	sid(h, w, r)
	noCache(header)

	s := h.pool.get(sessid)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	var data []byte
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	r.Body.Close()
	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		m, err := url.ParseQuery(buf.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data = []byte(m.Get("d"))
	case "text/plain":
		data = buf.Bytes()
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Payload expected."))
		return
	}

	var messages []string
	if err := json.Unmarshal(data, &messages); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Broken JSON encoding."))
		return
	}

	for _, v := range messages {
		s.in.push([]byte(v))
	}

	w.Write([]byte("ok"))
}
