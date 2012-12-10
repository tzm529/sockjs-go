package main

import (
	"github.com/fzzbt/sockjs-go/sockjs"
	"fmt"
	"net/http"
)

func echoHandler(s *sockjs.Session) {
	fmt.Println("session opened")
    for {
        m, err := s.Receive()
        if err != nil {
            break
        }
		fmt.Println("Received:", m)
        s.Send(m)
    }
	fmt.Println("session closing")
}

func main() {
	dwsconf := sockjs.NewConfig()
	dwsconf.Websocket = false

	http.Handle("/static", http.FileServer(http.Dir("./static")))
	sockjs.Handle("/echo", echoHandler, sockjs.NewConfig())
	sockjs.Handle("/disabled_websocket_echo", echoHandler, dwsconf)
	sockjs.Handle("/close", 
		func(s *sockjs.Session) { s.Close() }, 
		sockjs.NewConfig())
	err := http.ListenAndServe(":8081", sockjs.DefaultServeMux)
	if err != nil { fmt.Println(err) }
}
