package sockjs

import (
	"io"
)

type eventSourceProtocol struct{}

func (p eventSourceProtocol) contentType() string { return "text/event-stream; charset=UTF-8" }

func (p eventSourceProtocol) writePrelude(w io.Writer) (err error) {
	_, err = io.WriteString(w, "\r\n")
	return
}

func (p eventSourceProtocol) writeOpen(w io.Writer) (err error) {
	_, err = io.WriteString(w, "data: o\r\n\r\n")
	return
}

func (p eventSourceProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) {
	n, err = w.Write(aframe("data: ", "\r\n\r\n", m...))
	return
}

func (p eventSourceProtocol) writeClose(w io.Writer, code int, m string) {
	w.Write(cframe("data: ", code, m, "\r\n\r\n"))
}

func (p eventSourceProtocol) protocol() Protocol { return ProtocolEventSource }
