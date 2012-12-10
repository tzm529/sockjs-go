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

func receiveWebsocket(s *Session) (string, error) {
	var messages []string
	var data []byte

	err := websocket.Message.Receive(s.ws, &data)
	if err != nil {
		return "", err
	}

	// ignore empty frames
	if len(data) == 0 {
		return receiveWebsocket(s)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return "", err
	}

	// ignore empty messages
	if len(messages) == 0 {
		return receiveWebsocket(s)
	}

	if len(messages) > 1 {
		// push the leftover messages to the queue
		for _, v  := range messages[1:] {
			s.push(v)
		}
	}

	return messages[0], nil
}

func sendWebsocket(s *Session, m string) (err error) {
	_, err = s.ws.Write(aframe(m))
	return
}

func closeWebsocket(s *Session) error {
	s.ws.Write([]byte(`c[3000,"Go away!"]`))
	return s.ws.Close()
}
