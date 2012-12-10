package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"net/http"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := websocket.Handler(func(ws *websocket.Conn) {
		// initiate connection
		_, err := ws.Write([]byte{'o'})
		if err != nil {
			return
		}

		session := new(Session)
		session.kind = sessionKindWebsocket
		session.ws = ws
		s.hfunc(session)
	})

	h.ServeHTTP(w, r)
}

func receiveWebsocket(ws *websocket.Conn) (string, error) {
	var messages []string
	var data []byte

	err := websocket.Message.Receive(ws, &data)
	if err != nil {
		return "", err
	}

	// ignore empty frames
	if len(data) == 0 {
		return receiveWebsocket(ws)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return "", err
	}

	// ignore empty messages
	if len(messages) == 0 {
		return receiveWebsocket(ws)
	}

	if len(messages) > 1 {
		println("YLIYKS")
	}

	return messages[0], nil
}

func sendWebsocket(ws *websocket.Conn, s string) (err error) {
	_, err = ws.Write(aframe(s))
	return
}

func closeWebsocket(ws *websocket.Conn) error {
	ws.Write([]byte(`c[3000,"Go away!"]`))
	return ws.Close()
}
