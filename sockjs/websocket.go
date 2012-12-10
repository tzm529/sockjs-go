package sockjs

import (
	"io"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"net/http"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := websocket.Handler(func(ws *websocket.Conn) {
		var messages []string
		var data []byte

		// initiate connection
		_, err := ws.Write([]byte{'o'})
		if err != nil {
			return
		}

		conn := new(Conn)
		conn.kind = connKindWebsocket
		conn.wc = ws
		if s.OnOpen != nil { s.OnOpen(conn) }

		// Read messages until we get some error, like connection closed.
		for {
			err = websocket.Message.Receive(ws, &data) 
			if err != nil {
				break
			}

			// ignore empty frames
			if len(data) == 0 {
					continue
			}

			err = json.Unmarshal(data, &messages)
			if err != nil {
				break
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

func sendWebsocket(w io.Writer, s string) (err error) {
	_, err = w.Write(aframe(s))
	return
}

func closeWebsocket(wc io.WriteCloser) error {
	wc.Write([]byte(`c[3000,"Go away!"]`))
	return wc.Close()
}


