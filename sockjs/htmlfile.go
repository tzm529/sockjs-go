package sockjs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	htmlFileFormat string = `<!doctype html>
<html><head>
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head><body><h2>Don't panic!</h2>
  <script>
    document.domain = document.domain;
    var c = parent.%s;
    c.start();
    function p(d) {c.message(d);};
    window.onload = function() {c.stop();};
  </script>`
)

type htmlfileProtocol struct {
	callback string
}

func (p *htmlfileProtocol) contentType() string { return "text/html; charset=UTF-8" }

func (p *htmlfileProtocol) write(w io.Writer, m []byte) (n int, err error) {
	js, _ := json.Marshal(string(m))
	n, err = fmt.Fprintf(w, "<script>\np(%s);\n</script>\r\n", js)
	return
}

func (p *htmlfileProtocol) protocol() Protocol       { return ProtocolHtmlfile }
func (p *htmlfileProtocol) streaming() preludeWriter { return p }

func (p *htmlfileProtocol) writePrelude(w io.Writer) (err error) {
	prelude := fmt.Sprintf(htmlFileFormat, p.callback)
	if len(prelude) < 1024 {
		prelude += strings.Repeat(" ", 1024)
	}
	prelude += "\r\n"
	_, err = io.WriteString(w, prelude)
	return
}

func htmlfileHandler(h *Handler, w http.ResponseWriter, r *http.Request, sessid string) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	p := new(htmlfileProtocol)
	p.callback = callback
	protocolHandler(h, w, r, sessid, p)
}
