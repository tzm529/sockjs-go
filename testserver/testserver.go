package main

import (
	"fmt"
	"github.com/fzzbt/sockjs-go/sockjs"
	"net/http"
)

func echoHandler(s sockjs.Session) {
	fmt.Println("session opened")
	
	for {
		m, err := s.Receive()
		if err != nil {
			println("ERR:",err.Error())
			break
		}
		fmt.Println("Received:", string(m))
		s.Send(m)
	}
	fmt.Println("session closing")
}

func main() {
	server := sockjs.NewServer(http.DefaultServeMux)
	conf := sockjs.NewConfig()
	conf.ResponseLimit = 4096
	dwsconf := conf
	dwsconf.Websocket = false
	cookieconf := conf
	cookieconf.Jsessionid = true

	http.Handle("/static", http.FileServer(http.Dir("./static")))
	server.Handle("/echo", echoHandler, conf)
	server.Handle("/disabled_websocket_echo", echoHandler, dwsconf)
	server.Handle("/cookie_needed_echo", echoHandler, cookieconf)
	server.Handle("/close",
		func(s sockjs.Session) { s.Close() },
		conf)
	err := http.ListenAndServe(":8081", server)
	if err != nil {
		fmt.Println(err)
	}
}
