package sockjs_test

import (
	"fmt"
	"github.com/fzzy/sockjs-go/sockjs"
	"net/http"
)

func echoHandler(s sockjs.Session) {
	for {
		m := s.Receive()
		if m == nil {
			break
		}
		s.Send(m)
	}
}

func Example() {
	// Handlers for two echo servers and a file server.

	conf := sockjs.NewConfig()
	dwsconf := sockjs.NewConfig()
	dwsconf.Websocket = false

	mux := sockjs.NewServeMux(http.DefaultServeMux)
	mux.Handle("/echo", echoHandler, conf)
	mux.Handle("/disabled_websocket_echo", echoHandler, dwsconf)
	http.Handle("/static", http.FileServer(http.Dir("./static")))

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNewHandler() {
	// Handle only SockJS requests prefixed with "/echo".
	h := sockjs.NewHandler("/echo", echoHandler, sockjs.NewConfig())
	err := http.ListenAndServe(":8081", h)
	if err != nil {
		fmt.Println(err)
	}
}
