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

func receiveRawWebsocket(s *Session) (string, error) {
	var data string
	err := websocket.Message.Receive(s.ws, &data)
	if err != nil {
		return "", err
	}
	return data, nil
}

func sendRawWebsocket(s *Session, m string) (err error) {
	_, err = s.ws.Write([]byte(m))
	return
}

func closeRawWebsocket(s *Session) error {
	return s.ws.Close()
}
