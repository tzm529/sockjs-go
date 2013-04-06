package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

type rawWebsocketSession struct {
	ws   *websocket.Conn
	info *RequestInfo
}

func (s *rawWebsocketSession) Receive() (data []byte) {
	err := websocket.Message.Receive(s.ws, &data)
	if err != nil {
		return nil
	}
	return data
}

func (s *rawWebsocketSession) Send(m []byte) {
	s.ws.Write(m)
}

func (s *rawWebsocketSession) End() {
	s.Close(3000, "Go away!")
}

func (s *rawWebsocketSession) Close(code int, reason string) {
	// BUG(fzzy): 
	// Websocket.Close() should specify code and reason.
	// Websocket package does not allow doing this.
	// http://code.google.com/p/go/issues/detail?id=4588
	s.ws.Close()
}

func (s *rawWebsocketSession) Info() RequestInfo  { return *s.info }
func (s *rawWebsocketSession) Protocol() Protocol { return ProtocolRawWebsocket }
func (s *rawWebsocketSession) String() string     { return s.Info().RemoteAddr }

func rawWebsocketHandler(h *handler, w http.ResponseWriter, r *http.Request) {
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
