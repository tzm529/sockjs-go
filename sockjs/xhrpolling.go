package sockjs

import (
	"io"
)

type xhrPollingProtocol struct{}

func (p xhrPollingProtocol) contentType() string { return "application/javascript; charset=UTF-8"}

func (p xhrPollingProtocol) writeOpen(w io.Writer) (err error) { 
	_, err = io.WriteString(w, "o\n")
	return
}

func (p xhrPollingProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) { 
	n, err = w.Write(frame("", "\n", m...))
	return
}

func (p xhrPollingProtocol) writeClose(w io.Writer, code int, m string) { 
	w.Write(cframe("", code, m, "\n"))
}

