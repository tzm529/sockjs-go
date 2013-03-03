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

type websocketClosure struct {
	abrupt bool
	code   int
	reason string
}

type websocketSession struct {
	// read-only
	config *Config
	info   *RequestInfo
	sessid string
	ws     *websocket.Conn

	closer   chan *websocketClosure
	hbTicker *time.Ticker

	// lock for making Receive() thread-safe
	rio  sync.Mutex
	rbuf [][]byte

	mu     sync.RWMutex
	closed bool

}

//* Public methods

func (s *websocketSession) Receive() []byte {
	s.rio.Lock()
	defer s.rio.Unlock()
	var m []byte

	if len(s.rbuf) > 0 {
		m, s.rbuf = s.rbuf[0], s.rbuf[1:]
		return m
	}

again:
	// Empty buffer, read some messages to it and return the first one.
	var data []byte
	var messages []string

	err := websocket.Message.Receive(s.ws, &data)
	if err != nil {
		goto disconnect
	}

	// ignore, no frame
	if len(data) == 0 {
		goto again
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		goto disconnect
	}

	// ignore, no messages
	if len(messages) == 0 {
		goto again
	}

	for _, v := range messages {
		s.rbuf = append(s.rbuf, []byte(v))
	}

	m, s.rbuf = s.rbuf[0], s.rbuf[1:]
	return m

disconnect:
	logPrintf(s.config.Logger, "%s: receive error: %s\n", s, err)
	s.abruptClose()
	return nil
}

func (s *websocketSession) Send(m []byte) {
	_, err := s.ws.Write(aframe(m))
	if err != nil {
		logPrintf(s.config.Logger, "%s: send error: %s\n", s, err)
	}
}

func (s *websocketSession) End() {
	s.Close(3000, "Go away!")
}

func (s *websocketSession) Close(code int, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	closure := new(websocketClosure)
	closure.abrupt = false
	closure.code = code
	closure.reason = reason

	s.closer <- closure
}

func (s *websocketSession) Info() RequestInfo  { return *s.info }
func (s *websocketSession) Protocol() Protocol { return ProtocolWebsocket }

// for logging purposes
func (s *websocketSession) String() string {
	return s.info.RemoteAddr + "/" + s.sessid
}

//* Private methods

func (s *websocketSession) backend() {
	logPrintf(s.config.Logger, "%s: session opened\n", s)
loop:
	for {
		select {
		case <-s.hbTicker.C:
			_, err := s.ws.Write([]byte{'h'})
			if err != nil {
				logPrintf(s.config.Logger, "%s: heartbeat error: %s\n", s, err)
				s.mu.Lock()
				s.closed = true
				s.mu.Unlock()
				break loop
			}

		case c := <-s.closer:
			if !c.abrupt {
				s.ws.Write(cframe(c.code, c.reason))
			}
			break loop
		}
	}

	s.hbTicker.Stop()
	s.ws.Close()
	logPrintf(s.config.Logger, "%s: session closed\n", s)
}

func (s *websocketSession) abruptClose() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	closure := new(websocketClosure)
	closure.abrupt = true

	s.closer <- closure
}

func websocketHandler(h *handler, w http.ResponseWriter, r *http.Request, sessid string) {
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
		s.config = h.config
		s.info = newRequestInfo(r, h.prefix, s.config.Headers)
		s.sessid = sessid
		s.closer = make(chan *websocketClosure)
		s.hbTicker = time.NewTicker(s.config.HeartbeatDelay)
		go s.backend()
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
