package main

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

func main() {
	mux := sockjs.NewServeMux(http.DefaultServeMux)
	conf := sockjs.NewConfig()
	conf.ResponseLimit = 4096
	dwsconf := conf
	dwsconf.Websocket = false
	cookieconf := conf
	cookieconf.Jsessionid = true

	http.Handle("/static", http.FileServer(http.Dir("./static")))
	mux.Handle("/echo", echoHandler, conf)
	mux.Handle("/disabled_websocket_echo", echoHandler, dwsconf)
	mux.Handle("/cookie_needed_echo", echoHandler, cookieconf)
	mux.Handle("/close",
		func(s sockjs.Session) { s.End() },
		conf)
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		fmt.Println(err)
	}
}
