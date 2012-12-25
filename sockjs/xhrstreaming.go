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

type xhrStreamingProtocol struct{}

func (p xhrStreamingProtocol) contentType() string { return "application/javascript; charset=UTF-8"}

func (p xhrStreamingProtocol) writePrelude(w io.Writer) (err error) { 
	_, err = w.Write(prelude)
	return
}

func (p xhrStreamingProtocol) writeOpen(w io.Writer) (err error) { 
	_, err = io.WriteString(w, "o\n")
	return
}

func (p xhrStreamingProtocol) writeData(w io.Writer, m ...[]byte) (n int, err error) { 
	n, err = w.Write(frame("", "\n", m...))
	return
}

func (p xhrStreamingProtocol) writeClose(w io.Writer, code int, m string) { 
	w.Write(cframe("", code, m, "\n"))
}
