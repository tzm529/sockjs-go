package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"net/http"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := websocket.Handler(func(ws *websocket.Conn) {
		var messages []string

		// initiate connection
		_, err := ws.Write([]byte{'o'})
		if err != nil {
			return
		}

		conn := newConn(ws)
		if s.OnOpen != nil { s.OnOpen(conn) }

		// Read messages until we get some error, like connection closed.
		for {
			if err = websocket.JSON.Receive(ws, &messages); err != nil {
				// ignore empty frames
				if jsonerr, ok := err.(*json.SyntaxError); ok && jsonerr.Offset == 0 {
					continue
				} else {
					break
				}
			}

			if s.OnMessage != nil {
				for _, message := range messages {
					s.OnMessage(conn, message)
				}
			}
		}

		ws.Close()
		if s.OnClose != nil { s.OnClose(conn) }
	})

	h.ServeHTTP(w, r)
}
