package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

type sessionRawWebsocket struct{ ws *websocket.Conn }

func (p *sessionRawWebsocket) Receive() (data []byte, err error) {
	err = websocket.Message.Receive(p.ws, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *sessionRawWebsocket) Send(m []byte) (err error) {
	_, err = p.ws.Write(m)
	return
}

func (p *sessionRawWebsocket) Close() error {
	return p.ws.Close()
}

func handleRawWebsocket(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	wh := websocket.Handler(func(ws *websocket.Conn) {
		s := new(sessionRawWebsocket)
		s.ws = ws
		h.hfunc(s)
	})

	wh.ServeHTTP(w, r)
}
