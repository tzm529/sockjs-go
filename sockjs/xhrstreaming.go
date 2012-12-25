package sockjs

import (
	"io"
)

var prelude []byte = make([]byte, 2049)
func init() {
	for i := 0; i < 2048; i++ {
		prelude[i] = byte('h')
	}
	prelude[2048] = byte('\n')
}


type protoXhrStreaming struct{}

func (p protoXhrStreaming) contentType() string { return "application/javascript; charset=UTF-8"}
func (p protoXhrStreaming) writePrelude(w io.Writer) (err error) { 
	_, err = w.Write(prelude)
	return
}
func (p protoXhrStreaming) writeOpen(w io.Writer) (err error) { 
	_, err = w.Write([]byte("o\n"))
	return
}
func (p protoXhrStreaming) writeData(w io.Writer, m ...[]byte) (n int, err error) { 
	n, err = w.Write(frame("", "\n", m...))
	return
}
func (p protoXhrStreaming) writeClose(w io.Writer, code int, m string) { 
	w.Write(cframe("", code, m, "\n"))
}
