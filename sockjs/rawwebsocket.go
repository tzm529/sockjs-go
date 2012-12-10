package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)


func handleRawWebsocket(w http.ResponseWriter, r *http.Request, s *Handler) {
	h := websocket.Handler(func(ws *websocket.Conn) {
		session := new(Session)
		session.kind = sessionKindRawWebsocket
		session.ws = ws
		s.hfunc(session)
	})

	h.ServeHTTP(w, r)
}

func receiveRawWebsocket(ws *websocket.Conn) (string, error) {
	var data string
	err := websocket.Message.Receive(ws, &data)
	if err != nil {
		return "", err
	}
 	return data, nil
}


func sendRawWebsocket(ws *websocket.Conn, s string) (err error) {
	_, err = ws.Write([]byte(s))
	return
}

func closeRawWebsocket(ws *websocket.Conn) error {
	ws.Write([]byte("Go away!"))
	return ws.Close()
}
