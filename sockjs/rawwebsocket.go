package sockjs

import (
	"io"
	"code.google.com/p/go.net/websocket"
	"net/http"
)

func handleRawWebsocket(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := websocket.Handler(func(ws *websocket.Conn) {
		var message string

		conn := new(Conn)
		conn.kind = connKindRawWebsocket
		conn.wc = ws
		if s.OnOpen != nil { s.OnOpen(conn) }

		// Read messages until we get some error, like connection closed.
		for {
			if err := websocket.Message.Receive(ws, &message); err != nil {
				break
			}

			if s.OnMessage != nil {
				s.OnMessage(conn, message)
			}
		}

		ws.Close()
		if s.OnClose != nil { s.OnClose(conn) }
	})

	h.ServeHTTP(w, r)
}

func sendRawWebsocket(w io.Writer, s string) (err error) {
	_, err = w.Write([]byte(s))
	return
}

func closeRawWebsocket(wc io.WriteCloser) error {
	wc.Write([]byte("Go away!"))
	return wc.Close()
}
