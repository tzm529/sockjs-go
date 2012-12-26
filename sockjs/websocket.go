package sockjs

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type websocketSession struct {
	in *queue
	ws *websocket.Conn
}

func (s *websocketSession) Receive() (m []byte, err error) {
	m, err = s.in.pullNow()
	if err != nil {
		return nil, ErrSessionClosed
	}
	if m != nil {
		return
	}

	//* read some messages to the queue and pull the first one
	var messages []string
	var data []byte

	err = websocket.Message.Receive(s.ws, &data)
	if err != nil {
		return nil, err
	}

	// ignore, no frame
	if len(data) == 0 {
		return s.Receive()
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, err
	}

	// ignore, no messages
	if len(messages) == 0 {
		return s.Receive()
	}

	for _, v := range messages {
		s.in.push([]byte(v))
	}

	m, err = s.in.pull()
	if err != nil {
		return nil, ErrSessionClosed
	}

	return m, nil
}

func (s *websocketSession) Send(m []byte) (err error) {
	_, err = s.ws.Write(aframe("", "", m))
	return
}

func (s *websocketSession) Close() (err error) {
	s.ws.Write(cframe("", 3000, "Go away!", ""))
	err = s.ws.Close()
	return
}

func websocketHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	if !h.config.Websocket {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("Sec-WebSocket-Version") == "13" && r.Header.Get("Origin") == "" {
		r.Header.Set("Origin", r.Header.Get("Sec-WebSocket-Origin"))
	}
	if strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		http.Error(w, `Can "Upgrade" only to "WebSocket".`, http.StatusBadRequest)
		return
	}

	conn := strings.ToLower(r.Header.Get("Connection"))
	if conn == "keep-alive, upgrade" {
		r.Header.Set("Connection", "Upgrade")
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
		s.in = newQueue()
		s.ws = ws
		h.hfunc(s)
	})

	wh.ServeHTTP(w, r)
}

func websocketPostHandler(w http.ResponseWriter, r *http.Request) {
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
	return
}
