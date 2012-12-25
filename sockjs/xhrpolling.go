package sockjs

import (
	"io"
)

type protoXhrPolling struct{}

func (p protoXhrPolling) contentType() string { return "application/javascript; charset=UTF-8"}

func (p protoXhrPolling) writeOpen(w io.Writer) (err error) { 
	_, err = w.Write([]byte("o\n"))
	return
}

func (p protoXhrPolling) writeData(w io.Writer, m ...[]byte) (n int, err error) { 
	n, err = w.Write(frame("", "\n", m...))
	return
}

func (p protoXhrPolling) writeClose(w io.Writer, code int, m string) { 
	w.Write(cframe("", code, m, "\n"))
}

