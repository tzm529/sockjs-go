package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

type protoRawWebsocket struct { ws *websocket.Conn }

func (p protoRawWebsocket) Receive() (data []byte, err error) {
	err = websocket.Message.Receive(p.ws, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p protoRawWebsocket) Send(m []byte) (err error) {
	_, err = p.ws.Write(m)
	return
}

func (p protoRawWebsocket) Close() error {
	return p.ws.Close()
}

func handleRawWebsocket(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	wh := websocket.Handler(func(ws *websocket.Conn) {
		session := protoRawWebsocket{ws}
		h.hfunc(session)
	})

	wh.ServeHTTP(w, r)
}

