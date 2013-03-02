package sockjs

import (
	"fmt"
	"io"
)

type eventSourceProtocol struct{}

func (p eventSourceProtocol) contentType() string { return "text/event-stream; charset=UTF-8" }

func (p eventSourceProtocol) write(w io.Writer, m []byte) (n int, err error) {
	n, err = fmt.Fprintf(w, "data: %s\r\n\r\n", m)
	return
}

func (p eventSourceProtocol) protocol() Protocol       { return ProtocolEventSource }
func (p eventSourceProtocol) streaming() preludeWriter { return p }

func (p eventSourceProtocol) writePrelude(w io.Writer) (err error) {
	_, err = io.WriteString(w, "\r\n")
	return
}
