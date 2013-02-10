package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type websocketSession struct {
	ws   *websocket.Conn
	info *RequestInfo
	heartbeatDelay time.Duration

	rio sync.Mutex
	rbuf [][]byte

	wio sync.Mutex
	closed bool
	closeErr error
}

func (s *websocketSession) Receive() ([]byte, error) {
	s.rio.Lock()
	defer s.rio.Unlock()

	if len(s.rbuf) > 0 {
		m := s.rbuf[0]
		s.rbuf = s.rbuf[1:]
		return m, nil
	}
	
	// nil the buffer, so the underlying array of the old slice gets GC'd
	s.rbuf = nil

again:
	// Empty buffer, read some messages to it and return the first one.
	var messages []string
	var data []byte
	var m []byte

	err := websocket.Message.Receive(s.ws, &data)
	if err != nil { goto disconnect }

	// ignore, no frame
	if len(data) == 0 {
		goto again 
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {	goto disconnect }

	// ignore, no messages
	if len(messages) == 0 {
		goto again
	}

	for _, v := range messages {
		s.rbuf = append(s.rbuf, []byte(v))
	}
	m = s.rbuf[0]
	s.rbuf = s.rbuf[1:]
	return m, nil
disconnect:
	s.wio.Lock()
	s.disconnect()
	s.wio.Unlock()

	return nil, ErrSessionClosed
}

func (s *websocketSession) Send(m []byte) (err error) {
	s.wio.Lock()
	defer s.wio.Unlock()

	if s.closed { return ErrSessionClosed }
	_, err = s.ws.Write(aframe(m))
	return
}

func (s *websocketSession) Close() error {
	s.wio.Lock()
	defer s.wio.Unlock()

	// it must be safe to call Close() multiple times
	if s.closed { return s.closeErr }

	_, s.closeErr = s.ws.Write(cframe(3000, "Go away!"))
	s.ws.Close()
	s.closed = true
	return s.closeErr
}

func (s *websocketSession) disconnect() {
	s.ws.Close()
	s.closed = true
	s.closeErr = ErrSessionClosed
}

func (s *websocketSession) heartbeater() {
	for {
		time.Sleep(s.heartbeatDelay)
		s.wio.Lock()
		if s.closed { return }
		_, err := s.ws.Write([]byte{'h'})
		if err != nil {
			s.disconnect()
			return
		}			
		s.wio.Unlock()
	}
}

func (p *websocketSession) Info() RequestInfo  { return *p.info }
func (p *websocketSession) Protocol() Protocol { return ProtocolWebsocket }

func websocketHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	header := r.Header
	if header.Get("Sec-WebSocket-Version") == "13" && header.Get("Origin") == "" {
		header.Set("Origin", header.Get("Sec-WebSocket-Origin"))
	}
	if strings.ToLower(header.Get("Upgrade")) != "websocket" {
		http.Error(w, `Can "Upgrade" only to "WebSocket".`, http.StatusBadRequest)
		return
	}

	conn := strings.ToLower(header.Get("Connection"))
	if conn == "keep-alive, upgrade" {
		header.Set("Connection", "Upgrade")
	} else if conn != "upgrade" {
		http.Error(w, `"Connection" must be "Upgrade".`, http.StatusBadRequest)
		return
	}

	wh := websocket.Handler(func(ws *websocket.Conn) {
		// initiate connection
		_, err := ws.Write([]byte{'o'})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s := new(websocketSession)
		s.ws = ws
		s.info = newRequestInfo(r, h.prefix, h.config.Headers)
		s.heartbeatDelay = h.config.HeartbeatDelay
		go s.heartbeater()
		h.hfunc(s)
	})

	wh.ServeHTTP(w, r)
}

func websocketPostHandler(w http.ResponseWriter, r *http.Request) {
	// normal http methods don't seem to allow writing response without Content-Type
	conn, bufrw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Fprintf(bufrw,
		"HTTP/1.1 %d %s\r\n",
		http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed))
	fmt.Fprint(bufrw, "Content-Length: 0\r\n")
	fmt.Fprint(bufrw, "Allow: GET\r\n")
	fmt.Fprint(bufrw, "\r\n")
	bufrw.Flush()
}
