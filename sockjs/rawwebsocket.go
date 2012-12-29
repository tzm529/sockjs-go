package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

type rawWebsocketSession struct{ 
	ws *websocket.Conn
	info *RequestInfo
}

func (p *rawWebsocketSession) Receive() (data []byte, err error) {
	err = websocket.Message.Receive(p.ws, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *rawWebsocketSession) Send(m []byte) (err error) {
	_, err = p.ws.Write(m)
	return
}

func (p *rawWebsocketSession) Close() error {
	// BUG: Should specify close reason "Go away!".
	//      websocket package does not allow doing this.
	return p.ws.Close()
}

func (p *rawWebsocketSession) Info() RequestInfo { return *p.info }
func (s *rawWebsocketSession) Protocol() Protocol { return ProtocolRawWebsocket }

func rawWebsocketHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	wh := websocket.Handler(func(ws *websocket.Conn) {
		s := new(rawWebsocketSession)
		s.ws = ws
		s.info = newRequestInfo(r, h.prefix, h.config.Headers)
		h.hfunc(s)
	})

	wh.ServeHTTP(w, r)
}
