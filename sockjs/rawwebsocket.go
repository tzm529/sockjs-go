package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

func handleRawWebsocket(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	wh := websocket.Handler(func(ws *websocket.Conn) {
		session := new(Session)
		session.proto = protocolRawWebsocket
		session.ws = ws
		h.hfunc(session)
	})

	wh.ServeHTTP(w, r)
}

func receiveRawWebsocket(s *Session) (data []byte, err error) {
	err = websocket.Message.Receive(s.ws, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func sendRawWebsocket(s *Session, m []byte) (err error) {
	_, err = s.ws.Write(m)
	return
}

func closeRawWebsocket(s *Session) error {
	return s.ws.Close()
}
